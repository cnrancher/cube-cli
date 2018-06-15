package docker

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	ref "github.com/docker/distribution/reference"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
)

func UseLocalOrPull(ctx context.Context, dClient *client.Client, hostname string, containerImage string, prsMap map[string]PrivateRegistry) error {
	logrus.Debugf("checking image [%s] on host [%s]", containerImage, hostname)
	imageExists, err := localImageExists(ctx, dClient, hostname, containerImage)
	if err != nil {
		return err
	}
	if imageExists {
		logrus.Debugf("no pull necessary, image [%s] exists on host [%s]", containerImage, hostname)
		return nil
	}
	logrus.Infof("pulling image [%s] on host [%s]", containerImage, hostname)
	if err := pullImage(ctx, dClient, hostname, containerImage, prsMap); err != nil {
		return err
	}
	logrus.Infof("successfully pulled image [%s] on host [%s]", containerImage, hostname)
	return nil
}

func localImageExists(ctx context.Context, dClient *client.Client, hostname string, containerImage string) (bool, error) {
	logrus.Debugf("checking if image [%s] exists on host [%s]", containerImage, hostname)
	_, _, err := dClient.ImageInspectWithRaw(ctx, containerImage)
	if err != nil {
		if client.IsErrNotFound(err) {
			logrus.Debugf("image [%s] does not exist on host [%s]: %v", containerImage, hostname, err)
			return false, nil
		}
		return false, fmt.Errorf("error checking if image [%s] exists on host [%s]: %v", containerImage, hostname, err)
	}
	logrus.Debugf("image [%s] exists on host [%s]", containerImage, hostname)
	return true, nil
}

func pullImage(ctx context.Context, dClient *client.Client, hostname string, containerImage string, prsMap map[string]PrivateRegistry) error {
	pullOptions := types.ImagePullOptions{}
	regAuth, prURL, err := GetImageRegistryConfig(containerImage, prsMap)
	if err != nil {
		return err
	}
	if regAuth != "" && prURL == EngineRegistryURL {
		pullOptions.PrivilegeFunc = tryRegistryAuth(prsMap[prURL])
	}
	pullOptions.RegistryAuth = regAuth

	out, err := dClient.ImagePull(ctx, containerImage, pullOptions)
	if err != nil {
		return fmt.Errorf("can't pull docker image [%s] for host [%s]: %v", containerImage, hostname, err)
	}
	defer out.Close()
	if logrus.GetLevel() == logrus.DebugLevel {
		io.Copy(os.Stdout, out)
	} else {
		io.Copy(ioutil.Discard, out)
	}

	return nil
}

func GetImageRegistryConfig(image string, prsMap map[string]PrivateRegistry) (string, string, error) {
	namedImage, err := ref.ParseNormalizedNamed(image)
	if err != nil {
		return "", "", err
	}
	regURL := ref.Domain(namedImage)
	if pr, ok := prsMap[regURL]; ok {
		// We do this if we have some docker.io login information
		regAuth, err := getRegistryAuth(pr)
		return regAuth, pr.URL, err
	}
	return "", "", nil
}

func getRegistryAuth(pr PrivateRegistry) (string, error) {
	authConfig := types.AuthConfig{
		Username: pr.User,
		Password: pr.Password,
	}
	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(encodedJSON), nil
}

func tryRegistryAuth(pr PrivateRegistry) types.RequestPrivilegeFunc {
	return func() (string, error) {
		return getRegistryAuth(pr)
	}
}
