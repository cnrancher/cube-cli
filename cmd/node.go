package cmd

import (
	"fmt"

	"github.com/urfave/cli"
)

const (
	NodeDescription = `
Management Rancher Kubernetes Engine Node. 
					
Example:
	# List the Rancher Kubernetes Engine Nodes
	$ cube node ls
	# Add the Rancher Kubernetes Engine Node
	$ cube node add <node>
	# Remove the Rancher Kubernetes Engine Node
	$ cube node rm <node>
`
	Node = "node"
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
						Name:  Node,
						Usage: "Specify node name",
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
						Name:  Node,
						Usage: "Specify node name",
					},
				},
				Action: defaultAction(nodeRm),
			},
		},
	}
}

func nodeLs(ctx *cli.Context) error {
	fmt.Println("node list")
	return nil
}

func nodeAdd(ctx *cli.Context) error {
	node := ctx.String(Node)
	if "" == node {
		return fmt.Errorf("cube node add: require %v", Node)
	}
	fmt.Println("node add")
	return nil
}

func nodeRm(ctx *cli.Context) error {
	node := ctx.String(Node)
	if "" == node {
		return fmt.Errorf("cube node add: require %v", Node)
	}
	fmt.Println("node remove")
	return nil
}
