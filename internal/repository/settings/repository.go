package settings

import (
	"github.com/google/uuid"
	"github.com/yanodincov/k8s-forwarder/internal/config"
	"github.com/yanodincov/k8s-forwarder/internal/storage/yaml"
	"github.com/yanodincov/k8s-forwarder/pkg/helper"
)

const settingsStoragePath = "/settings.yaml"

type Repository struct {
	storage *yaml.Storage[Settings]
}

func NewStorage(cfg *config.Config) (*yaml.Storage[Settings], error) {
	return yaml.NewStorage[Settings](settingsStoragePath, cfg.AppName)
}

func NewRepository(storage *yaml.Storage[Settings]) *Repository {
	return &Repository{storage: storage}
}

func (r *Repository) GetSettings() (*Settings, error) {
	return r.storage.Get()
}

func (r *Repository) GetNamespaces() ([]NamespaceSetting, error) {
	settings, err := r.storage.Get()
	if err != nil {
		return nil, err
	}

	return settings.Namespaces, nil
}

func (r *Repository) GetNamespace(id uuid.UUID) (*NamespaceSetting, error) {
	settings, err := r.storage.Get()
	if err != nil {
		return nil, err
	}

	return settings.GetNamespace(id), nil
}

func (r *Repository) GetNamespaceByNameAndFile(namespace, configFilePath string) (*NamespaceSetting, error) {
	namespaces, err := r.GetNamespaces()
	if err != nil {
		return nil, err
	}

	res, ok := helper.SliceFind(namespaces, func(model NamespaceSetting) bool {
		return model.Namespace == namespace && model.ConfigFilePath == configFilePath
	})
	if !ok {
		return nil, nil
	}

	return &res, nil
}

func (r *Repository) AddNamespace(namespace NamespaceSetting) error {
	settings, err := r.storage.Get()
	if err != nil {
		return err
	}

	settings.AddNamespace(namespace)

	return r.storage.Save(settings)
}

func (r *Repository) RemoveNamespace(id uuid.UUID) error {
	settings, err := r.storage.Get()
	if err != nil {
		return err
	}

	settings.RemoveNamespace(id)
	if err = r.storage.Save(settings); err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetConfigFiles() ([]ConfigFileSetting, error) {
	settings, err := r.storage.Get()
	if err != nil {
		return nil, err
	}

	return settings.ConfigFiles, nil
}

func (r *Repository) GetConfigFile(id uuid.UUID) (*ConfigFileSetting, error) {
	settings, err := r.storage.Get()
	if err != nil {
		return nil, err
	}

	for _, configFile := range settings.ConfigFiles {
		if configFile.ID == id {
			return &configFile, nil
		}
	}

	return nil, nil
}

func (r *Repository) GetConfigFileByPath(path string) (*ConfigFileSetting, error) {
	settings, err := r.storage.Get()
	if err != nil {
		return nil, err
	}

	for _, configFile := range settings.ConfigFiles {
		if configFile.Path == path {
			return &configFile, nil
		}
	}

	return nil, nil
}

func (r *Repository) AddConfigFile(configFile ConfigFileSetting) error {
	settings, err := r.storage.Get()
	if err != nil {
		return err
	}

	settings.AddConfigFile(configFile)

	return r.storage.Save(settings)
}

func (r *Repository) RemoveConfigFile(id uuid.UUID) error {
	settings, err := r.storage.Get()
	if err != nil {
		return err
	}

	settings.RemoveConfigFile(id)

	return r.storage.Save(settings)
}
