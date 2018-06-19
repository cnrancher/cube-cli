package cmd

import (
	"context"
	"fmt"

	"github.com/cnrancher/cube-cli/docker"
	"github.com/cnrancher/cube-cli/util"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/go-connections/nat"
	"github.com/urfave/cli"
)

const (
	ServerDescription = `
Management RancherCUBE API-SERVER.

Example:
	# Run the RancherCUBE api-server
	$ cube server run --port "9600" --kube-config /example/kube-config.yml
	# Stop the RancherCUBE api-server
	$ cube server stop
	# Remove the RancherCUBE api-server
	$ cube server rm
	# Get the RancherCUBE api-server status
	$ cube server status
`

	ServerPort     = "port"
	ConfigLocation = "kube-config"
)

func ServerCommand() cli.Command {
	return cli.Command{
		Name:        "server",
		Aliases:     []string{"s"},
		Usage:       "Operations with cube api-server",
		Description: ServerDescription,
		Action:      defaultAction(serverStatus),
		Subcommands: []cli.Command{
			{
				Name:        "run",
				Usage:       "Run the RancherCUBE api-server",
				Description: "Run the RancherCUBE api-server",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:  ServerPort,
						Value: APIServerPortDefault,
						Usage: "Specify api-server listen port",
					},
					cli.StringFlag{
						Name:  ConfigLocation,
						Value: KubeConfigLocation,
						Usage: "Specify api-server kubernetes config location",
					},
				},
				Action: defaultAction(serverRun),
			},
			{
				Name:        "stop",
				Usage:       "Stop the RancherCUBE api-server",
				Description: "Stop the RancherCUBE api-server",
				Action:      defaultAction(serverStop),
			},
			{
				Name:        "rm",
				Usage:       "Remove the RancherCUBE api-server",
				Description: "Remove the RancherCUBE api-server",
				Action:      defaultAction(serverRm),
			},
			{
				Name:        "status",
				Usage:       "Status the RancherCUBE api-server",
				Description: "Status the RancherCUBE api-server",
				Action:      defaultAction(serverStatus),
			},
		},
	}
}

func serverRun(ctx *cli.Context) error {
	port := ctx.String(ServerPort)
	configLocation := ctx.String(ConfigLocation)
	if "" == configLocation {
		return fmt.Errorf("cube server run: require %v", ConfigLocation)
	}

	// generate JWT RSA256 key file
	err := util.GenerateRSA256()
	if err != nil {
		return err
	}

	context := context.Background()

	dClient, err := docker.NewClient(context, docker.SystemDockerSock)
	if err != nil {
		return err
	}

	// assemble *container.Config
	exposedPort, err := nat.NewPort("tcp", "9500")
	if err != nil {
		return err
	}
	exports := make(nat.PortSet, 1)
	exports[exposedPort] = struct{}{}
	containerConfig := &container.Config{
		Image: APIServerImage,
		Cmd: strslice.StrSlice{
			"serve",
			"--listen-addr=0.0.0.0:9500",
			"--kube-config=" + APIServerKubeConfig,
		},
		ExposedPorts: exports,
	}

	// assemble *container.HostConfig
	hostConfig := &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: util.RsaDirectory,
				Target: util.RsaDirectory,
			},
			{
				Type:   mount.TypeBind,
				Source: configLocation,
				Target: APIServerKubeConfig,
			},
		},
		PortBindings: nat.PortMap{
			"9500/tcp": []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: port + "/tcp",
				},
			},
		},
	}

	return docker.CreateOrRestart(context, dClient, containerConfig, hostConfig, nil, APIServerContainerName)
}

func serverStop(ctx *cli.Context) error {
	context := context.Background()

	dClient, err := docker.NewClient(context, docker.SystemDockerSock)
	if err != nil {
		return err
	}

	return docker.StopContainer(context, dClient, APIServerContainerName)
}

func serverRm(ctx *cli.Context) error {
	context := context.Background()

	dClient, err := docker.NewClient(context, docker.SystemDockerSock)
	if err != nil {
		return err
	}

	return docker.RemoveContainer(context, dClient, APIServerContainerName)
}

func serverStatus(ctx *cli.Context) error {
	context := context.Background()

	dClient, err := docker.NewClient(context, docker.SystemDockerSock)
	if err != nil {
		return err
	}

	_, err = docker.StatusContainer(context, dClient, APIServerContainerName)
	if err != nil {
		return err
	}

	return nil
}
