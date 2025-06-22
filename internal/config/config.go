package config

import (
	"gopkg.in/yaml.v3"
	"os"
	"time"
)

type ClientConfig struct {
	URL           string        `yaml:"url"`
	Token         string        `yaml:"token"`
	Timeout       time.Duration `yaml:"timeout"`
	RetryAttempts int           `yaml:"retry_attempts"`
	RetryDelay    time.Duration `yaml:"retry_delay"`
}

type gRPCServerConfig struct {
	Port string `yaml:"gRPCport"`
}

type HTTPServerConfig struct {
	Port         string        `yaml:"port"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
	IdleTimeout  time.Duration `yaml:"idle_timeout"`
}

type RateLimitConfig struct {
	RequestsPerMinute int `yaml:"requests_per_minute"`
}

type Config struct {
	ProductClient ClientConfig     `yaml:"product_client"`
	CartServer    HTTPServerConfig `yaml:"cart_server"`
	LOMSServer    gRPCServerConfig `yaml:"loms_server"`
	RateLimit     RateLimitConfig  `yaml:"rate_limit"`
}

func LoadConfig(path string) (*Config, error) {
	//absPath, err := filepath.Abs(path)
	//fmt.Println(absPath)
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cfg Config
	err = yaml.NewDecoder(file).Decode(&cfg)
	if err != nil {
		return nil, err
	}

	// Устанавливаем значения по умолчанию
	if cfg.ProductClient.Timeout == 0 {
		cfg.ProductClient.Timeout = 30 * time.Second
	}
	if cfg.ProductClient.RetryAttempts == 0 {
		cfg.ProductClient.RetryAttempts = 3
	}
	if cfg.ProductClient.RetryDelay == 0 {
		cfg.ProductClient.RetryDelay = 1 * time.Second
	}
	if cfg.CartServer.ReadTimeout == 0 {
		cfg.CartServer.ReadTimeout = 30 * time.Second
	}
	if cfg.CartServer.WriteTimeout == 0 {
		cfg.CartServer.WriteTimeout = 30 * time.Second
	}
	if cfg.CartServer.IdleTimeout == 0 {
		cfg.CartServer.IdleTimeout = 60 * time.Second
	}
	if cfg.RateLimit.RequestsPerMinute == 0 {
		cfg.RateLimit.RequestsPerMinute = 100
	}

	return &cfg, nil
}
