package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/cnrancher/cube-cli/util"
	rkecmd "github.com/rancher/rke/cmd"

	"github.com/rancher/types/apis/management.cattle.io/v3"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

const (
	NodeDescription = `
Management Rancher Kubernetes Engine Node. 
					
Example:
	# List the Rancher Kubernetes Engine Nodes
	$ cube node ls
	# Add the Rancher Kubernetes Engine Node
	$ cube node add --address <address> --roles worker,etcd --user rancher --ssh-key-path /home/rancher/.ssh/id_rsa
	# Remove the Rancher Kubernetes Engine Node
	$ cube node rm --address <address>
`
	Address    = "address"
	Roles      = "roles"
	User       = "user"
	SSHKeyPath = "ssh-key-path"
)

func NodeCommand() cli.Command {
	return cli.Command{
		Name:        "node",
		Aliases:     []string{"n"},
		Usage:       "Management Rancher Kubernetes Engine Node",
		Description: NodeDescription,
		Action:      defaultAction(nodeLs),
		Subcommands: []cli.Command{
			{
				Name:        "ls",
				Usage:       "List the Rancher Kubernetes Engine Nodes",
				Description: "List the Rancher Kubernetes Engine Nodes",
				Action:      defaultAction(nodeLs),
			},
			{
				Name:        "add",
				Usage:       "Add the Rancher Kubernetes Engine Node",
				Description: "Add the Rancher Kubernetes Engine Node",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  Address,
						Usage: "Specify node address, e.g. x.x.x.x",
					},
					cli.StringFlag{
						Name:  Roles,
						Value: "controlplane,worker,etcd",
						Usage: "Specify node roles",
					},
					cli.StringFlag{
						Name:  User,
						Value: "rancher",
						Usage: "Specify node user",
					},
					cli.StringFlag{
						Name:  SSHKeyPath,
						Value: "/home/rancher/.ssh/id_rsa",
						Usage: "Specify node ssh key path",
					},
				},

				Action: defaultAction(nodeAdd),
			},
			{
				Name:        "rm",
				Usage:       "Remove the Rancher Kubernetes Engine Node",
				Description: "Remove the Rancher Kubernetes Engine Node",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  Address,
						Usage: "Specify node name, e.g. x.x.x.x",
					},
				},
				Action: defaultAction(nodeRm),
			},
		},
	}
}

func nodeLs(ctx *cli.Context) error {
	config := &v3.RancherKubernetesEngineConfig{}
	var err error

	if _, fErr := os.Stat(NodeConfigDefault); fErr != nil {
		config, err = util.ReadRKEConfig(RKEConfigDefault)
	} else {
		config, err = util.ReadRKEConfig(NodeConfigDefault)
	}

	if err != nil {
		logrus.Errorf("%v", err)
		return err
	}
	logrus.Infof("%v", config)
	return nil
}

func nodeAdd(ctx *cli.Context) error {
	address := ctx.String(Address)
	if "" == address {
		return fmt.Errorf("cube node add: require %v", Address)
	}

	roles := ctx.String(Roles)
	rolesList := strings.Split(roles, ",")

	user := ctx.String(User)
	sshKeyPath := ctx.String(SSHKeyPath)

	config := &v3.RancherKubernetesEngineConfig{}
	var err error

	if _, fErr := os.Stat(NodeConfigDefault); fErr != nil {
		config, err = util.ReadRKEConfig(RKEConfigDefault)
	} else {
		config, err = util.ReadRKEConfig(NodeConfigDefault)
	}

	if err != nil {
		logrus.Errorf("%v", err)
		return err
	}

	if config.Nodes != nil && len(config.Nodes) > 0 {
		for _, node := range config.Nodes {
			if node.Address == address {
				logrus.Warnf("cube node add: node already exist")
				return nil
			}
		}
	}

	config.Nodes = append(config.Nodes, v3.RKEConfigNode{
		Address:    address,
		Role:       rolesList,
		User:       user,
		SSHKeyPath: sshKeyPath,
	})

	context := context.Background()
	_, _, _, _, _, err = rkecmd.ClusterUp(context, config, nil, nil, nil, false, "", false, true)
	if err != nil {
		logrus.Errorf("cube node add: execute rke command error %v", err)
		return err
	}

	err = util.WriteRKEConfig(config, NodeConfigDefault)
	if err != nil {
		logrus.Errorf("cube node add: write rke config file error %v", err)
		return err
	}

	return err
}

func nodeRm(ctx *cli.Context) error {
	address := ctx.String(Address)
	if "" == address {
		return fmt.Errorf("cube node remove: require %v", Address)
	}

	config := &v3.RancherKubernetesEngineConfig{}
	var err error

	if _, fErr := os.Stat(NodeConfigDefault); fErr != nil {
		config, err = util.ReadRKEConfig(RKEConfigDefault)
	} else {
		config, err = util.ReadRKEConfig(NodeConfigDefault)
	}

	if err != nil {
		logrus.Errorf("%v", err)
		return err
	}

	foundIndex := -1
	for index, node := range config.Nodes {
		if node.Address == address {
			foundIndex = index
			break
		}
	}

	if config.Nodes == nil || len(config.Nodes) <= 0 {
		logrus.Warnf("cube node remove: no nodes in config file")
		return nil
	}

	if len(config.Nodes) == 1 {
		config.Nodes = []v3.RKEConfigNode{}
	} else {
		left := config.Nodes[0 : foundIndex-1]
		right := config.Nodes[foundIndex+1:]
		config.Nodes = util.MergeNodes(left, right)
	}

	context := context.Background()
	_, _, _, _, _, err = rkecmd.ClusterUp(context, config, nil, nil, nil, false, "", false, true)
	if err != nil {
		logrus.Errorf("cube node remove: execute rke command error %v", err)
		return err
	}

	err = util.WriteRKEConfig(config, NodeConfigDefault)
	if err != nil {
		logrus.Errorf("cube node remove: write rke config error %v", err)
		return err
	}

	return err
}
