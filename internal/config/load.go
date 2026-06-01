package config

import "github.com/kelseyhightower/envconfig"

func Load() (*Config, error) {
	cfg := New()
	if err := envconfig.Process("", cfg); err != nil {
		return nil, err
	}
	return cfg, nil

}
