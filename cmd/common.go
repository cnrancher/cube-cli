package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var (
	APIServerKubeConfig    = "/etc/rancher/cube/kube-config.yml"
	APIServerImage         = "cnrancher/cube-apiserver"
	APIServerContainerName = "cube-apiserver"
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
