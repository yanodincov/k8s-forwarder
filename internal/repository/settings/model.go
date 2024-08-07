package settings

import (
	"github.com/google/uuid"
	"github.com/yanodincov/k8s-forwarder/pkg/helper"
)

type NamespaceSetting struct {
	ID             uuid.UUID `yaml:"id"`
	Namespace      string    `yaml:"namespace"`
	ConfigFilePath string    `yaml:"config_file_path"`
}

type ConfigFileSetting struct {
	ID   uuid.UUID `yaml:"id"`
	Path string    `yaml:"path"`
}

type Settings struct {
	Namespaces  []NamespaceSetting  `yaml:"namespaces"`
	ConfigFiles []ConfigFileSetting `yaml:"config_files"`
}

func (s *Settings) RemoveNamespace(id uuid.UUID) {
	for i, namespace := range s.Namespaces {
		if namespace.ID == id {
			s.Namespaces = append(s.Namespaces[:i], s.Namespaces[i+1:]...)
			return
		}
	}
}

func (s *Settings) AddNamespace(namespace NamespaceSetting) {
	s.Namespaces = append(s.Namespaces, namespace)
}

func (s *Settings) UpdateNamespace(namespace NamespaceSetting) {
	for i, n := range s.Namespaces {
		if n.ID == namespace.ID {
			s.Namespaces[i] = namespace
			return
		}
	}
}

func (s *Settings) GetNamespace(id uuid.UUID) *NamespaceSetting {
	if elem, ok := helper.SliceFind(s.Namespaces, func(item NamespaceSetting) bool {
		return item.ID == id
	}); ok {
		return &elem
	}

	return nil
}

func (s *Settings) RemoveConfigFile(id uuid.UUID) {
	for i, configFile := range s.ConfigFiles {
		if configFile.ID == id {
			s.ConfigFiles = append(s.ConfigFiles[:i], s.ConfigFiles[i+1:]...)
			return
		}
	}
}

func (s *Settings) AddConfigFile(configFile ConfigFileSetting) {
	s.ConfigFiles = append(s.ConfigFiles, configFile)
}

func (s *Settings) UpdateConfigFile(configFile ConfigFileSetting) {
	for i, c := range s.ConfigFiles {
		if c.ID == configFile.ID {
			s.ConfigFiles[i] = configFile
			return
		}
	}
}

func (s *Settings) GetConfigFile(id uuid.UUID) *ConfigFileSetting {
	if elem, ok := helper.SliceFind(s.ConfigFiles, func(item ConfigFileSetting) bool {
		return item.ID == id
	}); ok {
		return &elem
	}

	return nil
}

func (s *Settings) GetNamespacesByConfigFilePath(path string) []NamespaceSetting {
	return helper.SliceFilter(s.Namespaces, func(item NamespaceSetting) bool {
		return item.ConfigFilePath == path
	})
}
