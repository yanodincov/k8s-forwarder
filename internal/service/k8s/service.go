package k8s

import (
	"context"
	"github.com/pkg/errors"
	"github.com/yanodincov/k8s-forwarder/internal/repository/settings"
	"github.com/yanodincov/k8s-forwarder/pkg/helper"
	k8sclient "github.com/yanodincov/k8s-forwarder/pkg/k8s-client"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"slices"
	"time"
)

const timeout = 5 * time.Second

type Service struct {
	settingsRepository *settings.Repository
}

func NewService(settingsRepository *settings.Repository) *Service {
	return &Service{settingsRepository: settingsRepository}
}

func (s *Service) GetConfigFilePathsWithDefault() ([]string, error) {
	files, err := s.settingsRepository.GetConfigFiles()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get config files")
	}

	filePaths := helper.SliceMap(files, func(input settings.ConfigFileSetting) string {
		return input.Path
	})

	paths := make([]string, 0, len(files)+1)
	defaultEnvConfigFilePath := k8sclient.GetDefaultK8sConfigPath()
	if defaultEnvConfigFilePath != "" && !slices.Contains(filePaths, defaultEnvConfigFilePath) {
		paths = append(paths, defaultEnvConfigFilePath)
	}

	paths = append(paths, defaultEnvConfigFilePath)

	return paths, nil
}

func (s *Service) GetNamespaces(cfgFilePath string) ([]string, error) {
	client, err := k8sclient.NewClientFromFilePath(cfgFilePath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create k8s client")
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	k8sNamespacesList, err := client.CoreV1().Namespaces().List(ctx, metaV1.ListOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get k8s namespaces list")
	}

	return helper.SliceMap(k8sNamespacesList.Items, func(input coreV1.Namespace) string {
		return input.Name
	}), nil
}

func (s *Service) GetServices(cfgFilePath, namespace string) ([]coreV1.Service, error) {
	client, err := k8sclient.NewClientFromFilePath(cfgFilePath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create k8s client")
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	k8sServicesList, err := client.CoreV1().Services(namespace).List(ctx, metaV1.ListOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get k8s services list")
	}

	return k8sServicesList.Items, nil
}

func (s *Service) GetPodsByService(cfgFilePath, namespace, serviceName string) ([]coreV1.Pod, error) {
	client, err := k8sclient.NewClientFromFilePath(cfgFilePath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create k8s client")
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	k8sPodsList, err := client.CoreV1().Pods(namespace).List(ctx, metaV1.ListOptions{
		LabelSelector: "app=" + serviceName,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to get k8s pods list")
	}

	return k8sPodsList.Items, nil
}
