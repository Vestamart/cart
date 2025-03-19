package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type ClientConfig struct {
	URL   string `yaml:"url"`
	Token string `yaml:"token"`
}
type ServerConfig struct {
	Port string `yaml:"port"`
}

type Config struct {
	ProductClient ClientConfig `yaml:"product_client"`
	CartServer    ServerConfig `yaml:"cart_server"`
	LOMSServer    ServerConfig `yaml:"loms_server"`
}

func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = yaml.NewDecoder(file).Decode(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
