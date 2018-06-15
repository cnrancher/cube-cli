package docker

import "time"

var (
	EngineAPIVersion     = "1.24"
	EngineRegistryURL    = "docker.io"
	EngineDefaultSock    = "unix:///var/run/docker.sock"
	EngineSupportVersion = []string{"1.11.x", "1.12.x", "1.13.x", "17.03.x", "17.12.x", "18.03.x"}

	ContainerDefaultTimeout = 10 * time.Second

	SystemDockerSock = "unix:///var/run/system-docker.sock"
)

type PrivateRegistry struct {
	// URL for the registry
	URL string `yaml:"url" json:"url,omitempty"`
	// User name for registry acces
	User string `yaml:"user" json:"user,omitempty"`
	// Password for registry access
	Password string `yaml:"password" json:"password,omitempty"`
}
