package k8sclient

import (
	"github.com/pkg/errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/utils/env"
)

const k8sConfigPathEnv = "KUBECONFIG"

func GetDefaultK8sConfigPath() string {
	return env.GetString(k8sConfigPathEnv, "")
}

func NewClientFromEnvFilePath() (*kubernetes.Clientset, error) {
	return NewClientFromFilePath(GetDefaultK8sConfigPath())
}

func NewClientFromFilePath(filePath string) (*kubernetes.Clientset, error) {
	config, err := clientcmd.BuildConfigFromFlags("", filePath)
	if err != nil {
		panic(err.Error())
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create k8s client")
	}

	return client, nil
}
