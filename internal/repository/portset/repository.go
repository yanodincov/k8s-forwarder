package portset

import (
	"github.com/google/uuid"
	"github.com/yanodincov/k8s-forwarder/internal/config"
	"github.com/yanodincov/k8s-forwarder/internal/storage/yaml"
	"github.com/yanodincov/k8s-forwarder/pkg/helper"
)

const portSetStoragePath = "/portsets.yaml"

type Repository struct {
	storage *yaml.Storage[ServiceSets]
}

func NewStorage(cfg *config.Config) (*yaml.Storage[ServiceSets], error) {
	return yaml.NewStorage[ServiceSets](portSetStoragePath, cfg.AppName)
}

func NewRepository(storage *yaml.Storage[ServiceSets]) *Repository {
	return &Repository{storage: storage}
}

func (r *Repository) GetServiceSets() (*ServiceSets, error) {
	return r.storage.Get()
}

func (r *Repository) GetServiceSet(id uuid.UUID) (*ServiceSet, error) {
	configs, err := r.storage.Get()
	if err != nil {
		return nil, err
	}

	for _, config := range configs.Sets {
		if config.ID == id {
			return &config, nil
		}
	}

	return nil, nil
}

func (r *Repository) GetServiceSetByName(name string) (*ServiceSet, error) {
	configs, err := r.storage.Get()
	if err != nil {
		return nil, err
	}

	portSet, ok := helper.SliceFind(configs.Sets, func(set ServiceSet) bool {
		return set.Name == name
	})

	return helper.If(ok, &portSet, nil), nil
}

func (r *Repository) AddSet(config *ServiceSet) error {
	configs, err := r.storage.Get()
	if err != nil {
		return err
	}

	configs.AddSet(*config)

	return r.storage.Save(configs)
}

func (r *Repository) UpdateSet(config *ServiceSet) error {
	configs, err := r.storage.Get()
	if err != nil {
		return err
	}

	for i, c := range configs.Sets {
		if c.ID == config.ID {
			configs.Sets[i] = *config
			break
		}
	}

	return r.storage.Save(configs)
}

func (r *Repository) RemoveSet(id uuid.UUID) error {
	configs, err := r.storage.Get()
	if err != nil {
		return err
	}

	configs.RemoveSet(id)

	return r.storage.Save(configs)
}

func (r *Repository) GetService(config *ServiceSet) (*ServiceForwardConfig, error) {
	for _, service := range config.Services {
		if service.ID == config.ID {
			return &service, nil
		}
	}

	return nil, nil
}

func (r *Repository) AddService(service ServiceForwardConfig) error {
	configs, err := r.storage.Get()
	if err != nil {
		return err
	}

	for i, config := range configs.Sets {
		if config.ID == service.SetID {
			configs.Sets[i].AddServiceForwardConfig(service)
			break
		}
	}

	return r.storage.Save(configs)
}

func (r *Repository) UpdateService(service ServiceForwardConfig) error {
	configs, err := r.storage.Get()
	if err != nil {
		return err
	}

	for _, config := range configs.Sets {
		config.UpdateServiceForwardConfig(service)
	}

	return r.storage.Save(configs)
}

func (r *Repository) RemoveService(serviceID uuid.UUID) error {
	configs, err := r.storage.Get()
	if err != nil {
		return err
	}

	for i := range configs.Sets {
		configs.Sets[i].RemoveServiceForwardConfig(serviceID)
	}

	return r.storage.Save(configs)
}

type ServicePortFilter struct {
	SetID       *uuid.UUID
	ConfigPath  *string
	Namespace   *string
	ServiceName *string
	ServicePort *int
	LocalPort   *int
}

func (r *Repository) GetOneServicePortByFilter(filter ServicePortFilter) (*ServiceForwardConfig, error) {
	configs, err := r.storage.Get()
	if err != nil {
		return nil, err
	}

	for _, config := range configs.Sets {
		for _, service := range config.Services {
			if filter.SetID != nil && config.ID != *filter.SetID {
				continue
			}
			if filter.ConfigPath != nil && service.ConfigFilePath != *filter.ConfigPath {
				continue
			}
			if filter.Namespace != nil && service.Namespace != *filter.Namespace {
				continue
			}
			if filter.ServiceName != nil && service.ServiceName != *filter.ServiceName {
				continue
			}
			if filter.ServicePort != nil && service.ServicePort != *filter.ServicePort {
				continue
			}
			if filter.LocalPort != nil && service.LocalPort != *filter.LocalPort {
				continue
			}

			return &service, nil
		}
	}

	return nil, nil
}
