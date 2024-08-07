package portset

import "github.com/google/uuid"

type ServiceForwardConfig struct {
	ID             uuid.UUID `yaml:"id"`
	SetID          uuid.UUID `yaml:"set_id"`
	ServiceName    string    `yaml:"name"`
	ServicePort    int       `yaml:"service_port"`
	LocalPort      int       `yaml:"local_port"`
	Namespace      string    `yaml:"namespace"`
	ConfigFilePath string    `yaml:"config_file_path"`
}

type ServiceSet struct {
	ID       uuid.UUID              `yaml:"id"`
	Name     string                 `yaml:"name"`
	Services []ServiceForwardConfig `yaml:"services,omitempty"`
}

func (n *ServiceSet) AddServiceForwardConfig(service ServiceForwardConfig) {
	n.Services = append(n.Services, service)
}

func (n *ServiceSet) RemoveServiceForwardConfig(id uuid.UUID) {
	for i, service := range n.Services {
		if service.ID == id {
			n.Services = append(n.Services[:i], n.Services[i+1:]...)
			return
		}
	}
}

func (n *ServiceSet) UpdateServiceForwardConfig(service ServiceForwardConfig) {
	for i, s := range n.Services {
		if s.ID == service.ID {
			n.Services[i] = service
			return
		}
	}
}

type ServiceSets struct {
	Sets []ServiceSet `yaml:"sets,omitempty"`
}

func (c *ServiceSets) AddSet(config ServiceSet) {
	c.Sets = append(c.Sets, config)
}

func (c *ServiceSets) RemoveSet(id uuid.UUID) {
	for i, config := range c.Sets {
		if config.ID == id {
			c.Sets = append(c.Sets[:i], c.Sets[i+1:]...)
			return
		}
	}
}

func (c *ServiceSets) UpdateSet(config ServiceSet) {
	for i, set := range c.Sets {
		if set.ID == config.ID {
			c.Sets[i] = config
			return
		}
	}
}
