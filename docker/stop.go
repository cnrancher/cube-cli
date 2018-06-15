package docker

import (
	"context"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
)

func StopContainer(ctx context.Context, dClient *client.Client, containerName string) error {
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
		// stop the container
		if err := dClient.ContainerStop(ctx, containerID, &ContainerDefaultTimeout); err != nil {
			logrus.Errorf("stop container %s error: %v", containerID, err)
			return err
		}
		return nil
	}

	logrus.Warnf("container not found")

	return nil
}