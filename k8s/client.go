package k8s

import (
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	clientGenerator *ClientGenerator
)

type ClientGenerator struct {
	Clientset kubernetes.Clientset
}

func NewClientGenerator(kubeConfig string) *ClientGenerator {
	if clientGenerator == nil {
		var config *rest.Config
		var err error

		if kubeConfig == "" {
			config, err = rest.InClusterConfig()
			if err != nil {
				logrus.Fatalf("generate config failed: %v", err)
			}
		}

		config, err = clientcmd.BuildConfigFromFlags("", kubeConfig)
		if err != nil {
			logrus.Fatalf("generate config failed: %v", err)
		}

		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			logrus.Fatalf("generate clientset failed: %v", err)
		}

		clientGenerator = &ClientGenerator{
			Clientset: *clientset,
		}
	}

	return clientGenerator
}
