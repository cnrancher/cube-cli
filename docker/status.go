package docker

import (
	"context"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
)

func StatusContainer(ctx context.Context, dClient *client.Client, containerName string) (string, error) {
	containers, err := dClient.ContainerList(ctx, types.ContainerListOptions{All: true})
	if err != nil {
		logrus.Errorf("list docker containers error: %v", err)
		return "", err
	}

	isFound := false
	containerStatus := ""
	// check that the container with the specified name exists
	for _, c := range containers {
		for _, name := range c.Names {
			if strings.Contains(name, containerName) {
				isFound = true
				containerStatus = c.Status
				break
			}
		}
		if isFound {
			break
		}
	}

	if isFound && containerStatus != "" {
		// status the container
		logrus.Info(containerStatus)
		return containerStatus, nil
	}

	logrus.Warnf("container not found")

	return "NotFound", nil
}
