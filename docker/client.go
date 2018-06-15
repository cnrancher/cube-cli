package docker

import (
	"context"
	"strings"

	"github.com/coreos/go-semver/semver"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func NewClient(ctx context.Context, host string) (*client.Client, error) {
	if host == "" {
		return nil, errors.New("engine client host param is empty")
	}

	dClient, err := client.NewClient(host, EngineAPIVersion, nil, nil)
	if err != nil {
		return nil, err
	}

	err = CheckEngineVersion(ctx, dClient)
	if err != nil {
		return nil, err
	}

	return dClient, nil
}

func CheckEngineVersion(ctx context.Context, dClient *client.Client) error {
	info, err := dClient.Info(ctx)
	if err != nil {
		return errors.WithMessage(err, "can not retrieve engine info")
	}
	isValid, err := IsSupportVersion(info, EngineSupportVersion)
	if err != nil {
		return errors.Errorf("error while determining supported engine version [%s]: %v", info.ServerVersion, err)
	}
	if !isValid {
		logrus.Warnf("unsupported engine version found [%s], supported versions are %v", info.ServerVersion, EngineSupportVersion)
	}
	return nil
}

func IsSupportVersion(info types.Info, versions []string) (bool, error) {
	if strings.Contains(info.ServerVersion, "ros") {
		return true, nil
	}
	engineVersion, err := semver.NewVersion(info.ServerVersion)
	if err != nil {
		return false, err
	}
	for _, version := range versions {
		supportedEngineVersion, err := convertToSemver(version)
		if err != nil {
			return false, err
		}
		if engineVersion.Major == supportedEngineVersion.Major && engineVersion.Minor == supportedEngineVersion.Minor {
			return true, nil
		}
	}
	return false, nil
}

func convertToSemver(version string) (*semver.Version, error) {
	compVersion := strings.SplitN(version, ".", 3)
	if len(compVersion) != 3 {
		return nil, errors.New("the default version is not correct")
	}
	compVersion[2] = "0"
	return semver.NewVersion(strings.Join(compVersion, "."))
}
