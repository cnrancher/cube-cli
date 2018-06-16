package util

import (
	"io/ioutil"

	"github.com/rancher/types/apis/management.cattle.io/v3"
	"gopkg.in/yaml.v2"
)

func ReadRKEConfig(filename string) (*v3.RancherKubernetesEngineConfig, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	cluster := v3.RancherKubernetesEngineConfig{}
	err = yaml.Unmarshal(bytes, &cluster)
	if err != nil {
		return nil, err
	}

	return &cluster, nil
}

func WriteRKEConfig(config *v3.RancherKubernetesEngineConfig, filename string) error {
	bytes, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, bytes, 0640)
}

func MergeNodes(s ...[]v3.RKEConfigNode) (slice []v3.RKEConfigNode) {
	switch len(s) {
	case 0:
		break
	case 1:
		slice = s[0]
		break
	default:
		s1 := s[0]
		s2 := MergeNodes(s[1:]...)
		slice = make([]v3.RKEConfigNode, len(s1)+len(s2))
		copy(slice, s1)
		copy(slice[len(s1):], s2)
		break
	}

	return
}
