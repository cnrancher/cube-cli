package cmd

import (
	rkecmd "github.com/rancher/rke/cmd"
	"github.com/urfave/cli"
)

func RKECommand() cli.Command {
	return cli.Command{
		Name:        "rke",
		Usage:       "Mapping the RKE commands",
		Description: "Manage the RancherCUBE Kubernetes",
		Subcommands: []cli.Command{
			rkecmd.UpCommand(),
			rkecmd.RemoveCommand(),
			rkecmd.VersionCommand(),
			rkecmd.ConfigCommand(),
			rkecmd.EtcdCommand(),
		},
	}
}
