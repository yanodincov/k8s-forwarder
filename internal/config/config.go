package config

const appName = "k8s-forwarder"

type Config struct {
	AppName string
}

func NewConfig() *Config {
	return &Config{
		AppName: appName,
	}
}
