package cmd

import (
	"fmt"
	"os"
	"os/user"
	"strings"

	"github.com/cnrancher/cube-cli/cmd/pkg/table"
	"github.com/cnrancher/cube-cli/util"

	"github.com/cnrancher/cube-cli/k8s"
	"github.com/rancher/types/apis/management.cattle.io/v3"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"k8s.io/api/core/v1"
)

const (
	NodeDescription = `
Management Rancher Kubernetes Engine Node. 
					
Example:
	# List the Rancher Kubernetes Engine Nodes
	$ cube node ls
	# Add the Rancher Kubernetes Engine Node
	$ cube node add <address> --roles worker,etcd --user rancher --ssh-key-path <user_directory>/.ssh/id_rsa
	# Remove the Rancher Kubernetes Engine Node
	$ cube node rm <address>
`
	Address    = "address"
	Roles      = "roles"
	User       = "user"
	SSHKeyPath = "ssh-key-path"
)

type NodeOutput struct {
	Config v3.RKEConfigNode `yaml:"config,omitempty" json:"config,omitempty"`
	Sync   bool             `yaml:"sync,omitempty" json:"sync,omitempty"`
}

func assembleSSHKeyPath() string {
	home := "/home/rancher"
	user, err := user.Current()
	if err == nil {
		home = user.HomeDir
	}

	return fmt.Sprintf(SSHKeyPathDefault, home)
}

func NodeCommand() cli.Command {
	return cli.Command{
		Name:        "node",
		Aliases:     []string{"n"},
		Usage:       "Management Rancher Kubernetes Engine Node",
		Description: NodeDescription,
		Flags:       table.WriterNodeFlags(),
		Action:      defaultAction(nodeLs),
		Subcommands: []cli.Command{
			{
				Name:        "ls",
				Usage:       "List the Rancher Kubernetes Engine Nodes",
				Description: "List the Rancher Kubernetes Engine Nodes",
				Flags:       table.WriterNodeFlags(),
				Action:      defaultAction(nodeLs),
			},
			{
				Name:        "add",
				Usage:       "Add the Rancher Kubernetes Engine Node",
				Description: "Add the Rancher Kubernetes Engine Node",
				ArgsUsage:   "<address>",
				Flags: []cli.Flag{
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
						Value: assembleSSHKeyPath(),
						Usage: "Specify node ssh key path",
					},
				},

				Action: defaultAction(nodeAdd),
			},
			{
				Name:        "rm",
				Usage:       "Remove the Rancher Kubernetes Engine Node",
				Description: "Remove the Rancher Kubernetes Engine Node",
				ArgsUsage:   "<address>",
				Action:      defaultAction(nodeRm),
			},
		},
	}
}

func nodeLs(ctx *cli.Context) error {
	config := &v3.RancherKubernetesEngineConfig{}
	var err error

	if _, fErr := os.Stat(RKEConfigDefault); fErr != nil {
		config, err = util.ReadRKEConfig(RKEBaseConfigDefault)
	} else {
		config, err = util.ReadRKEConfig(RKEConfigDefault)
	}

	if err != nil {
		logrus.Errorf("%v", err)
		return err
	}

	// retrieve form kubernetes backend
	client := k8s.NewClientGenerator(KubeConfigLocation)
	nodes, err := client.Clientset.CoreV1().Nodes().List(util.ListEverything)
	if err != nil {
		logrus.Errorf("can not retrieve kubernetes nodes: %v", err)
		return err
	}

	// whether rke config file nodes match kubernetes nodes or not
	isSyncMap := map[string]bool{}
	for _, node := range nodes.Items {
		for _, address := range node.Status.Addresses {
			if address.Type == v1.NodeInternalIP {
				isSyncMap[address.Address] = true
				break
			}
		}
	}

	writer := table.NewNodeWriter([][]string{
		{"ADDRESS", "{{.Config.Address}}"},
		{"ROLE", "{{.Config.Role}}"},
		{"USER", "{{.Config.User}}"},
		{"SSH KEY PATH", "{{.Config.SSHKeyPath}}"},
		{"SYNC", "{{.Sync}}"},
	}, ctx)
	defer writer.Close()

	for _, node := range config.Nodes {
		if _, ok := isSyncMap[node.Address]; ok {
			output := &NodeOutput{
				Config: node,
				Sync:   true,
			}
			writer.Write(*output)
		} else {
			output := &NodeOutput{
				Config: node,
				Sync:   false,
			}
			writer.Write(*output)
		}
	}

	return writer.Err()
}

func nodeAdd(ctx *cli.Context) error {
	args := ctx.Args()
	if len(args) < 1 {
		return fmt.Errorf("cube node add: no arguments")
	}

	address := args[0]
	if "" == address {
		return fmt.Errorf("cube node add: require %v", Address)
	}

	roles := ctx.String(Roles)
	rolesList := strings.Split(roles, ",")

	user := ctx.String(User)
	sshKeyPath := ctx.String(SSHKeyPath)

	config := &v3.RancherKubernetesEngineConfig{}
	var err error

	if _, fErr := os.Stat(RKEConfigDefault); fErr != nil {
		config, err = util.ReadRKEConfig(RKEBaseConfigDefault)
	} else {
		config, err = util.ReadRKEConfig(RKEConfigDefault)
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

	err = util.WriteRKEConfig(config, RKEConfigDefault)
	if err != nil {
		logrus.Errorf("cube node add: write rke config file error %v", err)
		return err
	}

	return err
}

func nodeRm(ctx *cli.Context) error {
	args := ctx.Args()
	if len(args) < 1 {
		return fmt.Errorf("cube node remove: no arguments")
	}

	address := args[0]
	if "" == address {
		return fmt.Errorf("cube node remove: require %v", Address)
	}

	config := &v3.RancherKubernetesEngineConfig{}
	var err error

	if _, fErr := os.Stat(RKEConfigDefault); fErr != nil {
		config, err = util.ReadRKEConfig(RKEBaseConfigDefault)
	} else {
		config, err = util.ReadRKEConfig(RKEConfigDefault)
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
		left := config.Nodes[:foundIndex]
		right := config.Nodes[foundIndex+1:]
		config.Nodes = util.MergeNodes(left, right)
	}

	err = util.WriteRKEConfig(config, RKEConfigDefault)
	if err != nil {
		logrus.Errorf("cube node remove: write rke config error %v", err)
		return err
	}

	return err
}
