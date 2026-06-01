package config

import (
	"errors"
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

func Load() (*Config, error) {
	cfg := New()
	if err := envconfig.Process("", cfg); err != nil {
		return nil, err
	}
	if cfg.DBHost == "" {
		msg := "DBHost is empty"
		return nil, fmt.Errorf("load config: %w", errors.New(msg))
	}
	if cfg.HTTPPort < 1 || cfg.HTTPPort > 65535 {
		msg := "HTTPPort is invalid"
		return nil, fmt.Errorf("load config: %w", errors.New(msg))
	}
	return cfg, nil

}
