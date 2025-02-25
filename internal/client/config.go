package client

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
	Client ClientConfig `yaml:"client"`
	Server ServerConfig `yaml:"server"`
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
