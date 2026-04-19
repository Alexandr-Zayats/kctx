package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Aliases  map[string]string `yaml:"aliases"`
	Defaults struct {
		Namespace string `yaml:"namespace"`
	} `yaml:"defaults"`
}

func Load() (*Config, error) {
	home, _ := os.UserHomeDir()
	path := filepath.Join(home, ".kctx", "config.yaml")

	data, err := os.ReadFile(path)
	if err != nil {
		return &Config{}, nil // не ошибка
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
