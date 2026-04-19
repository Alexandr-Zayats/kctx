package aws

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Project string `yaml:"project"`

	SSO struct {
		StartURL string `yaml:"start_url"`
		Region   string `yaml:"region"`
		RoleName string `yaml:"role_name"`
	} `yaml:"sso"`

	Accounts map[string]struct {
		AccountID string `yaml:"account_id"`
		Region    string `yaml:"region"`
	} `yaml:"accounts"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
