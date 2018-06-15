package cmd

import (
	"os"

	rkecmd "github.com/rancher/rke/cmd"
	"github.com/urfave/cli"
)

const (
	RKEConfigDefault = "/var/lib/rancher/cube/rke_config.yml"
)

func RKECommand() cli.Command {
	return cli.Command{
		Name:        "rke",
		Usage:       "Mapping the RKE commands",
		Description: "Manage the RancherCUBE Kubernetes",
		Before: func(c *cli.Context) error {
			if os.Getenv("RKE_CONFIG") == "" {
				os.Setenv("RKE_CONFIG", RKEConfigDefault)
			}
			return nil
		},
		Subcommands: []cli.Command{
			rkecmd.UpCommand(),
			rkecmd.RemoveCommand(),
			rkecmd.VersionCommand(),
			rkecmd.EtcdCommand(),
		},
	}
}
