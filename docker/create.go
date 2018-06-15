package docker

import (
	"context"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
)

func CreateOrRestart(ctx context.Context, dClient *client.Client, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, containerName string) error {
	containers, err := dClient.ContainerList(ctx, types.ContainerListOptions{All: true})
	if err != nil {
		logrus.Errorf("list docker containers error: %v", err)
		return err
	}

	isFound := false
	containerID := ""
	// check that the container with the specified name exists
	for _, c := range containers {
		for _, name := range c.Names {
			if strings.ContainsAny(name, containerName) {
				isFound = true
				containerID = c.ID
				break
			}
		}
		if isFound {
			break
		}
	}

	if isFound && containerID != "" {
		// restart the container
		if err := dClient.ContainerRestart(ctx, containerID, &ContainerDefaultTimeout); err != nil {
			logrus.Errorf("restart container %s error: %v", containerID, err)
			return err
		}
		return nil
	}

	// pull container image
	hostName, err := os.Hostname()
	if err != nil {
		logrus.Errorf("get host name error: %v", err)
		return err
	}
	// TODO: currently only support docker.io, this need to be enhanced to support multiple registry
	parseMap := make(map[string]PrivateRegistry)
	for _, pr := range parseMap {
		if pr.URL == "" {
			pr.URL = EngineRegistryURL
		}
		parseMap[pr.URL] = pr
	}
	err = UseLocalOrPull(ctx, dClient, hostName, config.Image, parseMap)
	if err != nil {
		logrus.Errorf("use local or pull image %s error: %v", config.Image, err)
		return err
	}

	// create container
	resp, err := dClient.ContainerCreate(ctx, config, hostConfig, networkingConfig, containerName)
	if err != nil {
		logrus.Errorf("create container %s error: %v", containerID, err)
		return err
	}

	if err := dClient.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		logrus.Errorf("start container %s error: %v", containerID, err)
		return err
	}

	return nil
}
