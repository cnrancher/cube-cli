package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/cnrancher/cube-cli/cmd/pkg/table"
	"github.com/cnrancher/cube-cli/docker"

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
	$ cube server run --port "9600"
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
		Flags:       table.WriterServerFlags(),
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
				Flags:       table.WriterServerFlags(),
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

	if configLocation != KubeConfigLocation {
		err := os.Rename(configLocation, KubeConfigLocation)
		if err != nil {
			return err
		}
	}

	context := context.Background()

	dClient, err := docker.NewClient(context, docker.SystemDockerSock)
	if err != nil {
		return err
	}

	// assemble *container.Config
	exposedPort, err := nat.NewPort("tcp", "9600")
	if err != nil {
		return err
	}
	exports := make(nat.PortSet, 1)
	exports[exposedPort] = struct{}{}
	containerConfig := &container.Config{
		Image: APIServerImage,
		Cmd: strslice.StrSlice{
			"serve",
			"--listen-addr=0.0.0.0:9600",
		},
		ExposedPorts: exports,
	}

	// assemble *container.HostConfig
	hostConfig := &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: APIServerKubeConfig,
				Target: APIServerKubeConfig,
			},
		},
		PortBindings: nat.PortMap{
			"9600/tcp": []nat.PortBinding{
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

	container, err := docker.StatusContainer(context, dClient, APIServerContainerName)
	if err != nil {
		return err
	}

	writer := table.NewServerWriter([][]string{
		{"CONTAINER ID", "{{.ID | id}}"},
		{"IMAGE", "{{.Image}}"},
		{"COMMAND", "{{.Command | cmd}}"},
		{"CREATED", "{{.Created | ago}}"},
		{"STATUS", "{{.Status}}"},
		{"PORTS", "{{.Ports | port}}"},
		{"NAMES", "{{.Names | name}}"},
	}, ctx)
	defer writer.Close()

	if container.ID != "" {
		writer.Write(container)
	}

	return writer.Err()
}
