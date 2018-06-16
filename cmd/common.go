package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

const (
	APIServerKubeConfig    = "/var/lib/rancher/cube/kube-config.yml"
	APIServerImage         = "cnrancher/cube-apiserver"
	APIServerContainerName = "cube-apiserver"
	RKEConfigDefault       = "/var/lib/rancher/cube/rke_config.yml"
	NodeConfigDefault      = "/var/lib/rancher/cube/node_config.yml"
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
