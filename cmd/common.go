package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

const (
	APIServerKubeConfig    = "/var/lib/rancher/cube/kube-config.yml"
	APIServerImage         = "cnrancher/cube-apiserver"
	APIServerContainerName = "cube-apiserver"
	APIServerPortDefault   = "9600"
	RKEConfigDefault       = "/var/lib/rancher/cube/rke_base.yml"
	NodeConfigDefault      = "/var/lib/rancher/cube/node_config.yml"
	SSHKeyPathDefault      = "/home/rancher/.ssh/id_rsa"
)

func defaultAction(fn func(ctx *cli.Context) error) func(ctx *cli.Context) error {
	return func(ctx *cli.Context) error {
		if ctx.Bool("help") {
			cli.ShowAppHelp(ctx)
			return nil
		}

		if err := fn(ctx); err != nil {
			logrus.Errorf("cube error: %v", err)
			return err
		}
		return nil
	}
}
