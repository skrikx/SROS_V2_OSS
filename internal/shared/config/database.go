package config

import (
	"fmt"
	"net/url"
	"strings"
)

type DatabaseConfig struct {
	Enabled       bool   `json:"enabled"`
	Driver        string `json:"driver"`
	URL           string `json:"url"`
	MigrationsDir string `json:"migrations_dir"`
	Schema        string `json:"schema"`
}

func (d DatabaseConfig) Summary() map[string]any {
	return map[string]any{
		"enabled":        d.Enabled,
		"driver":         d.Driver,
		"url_present":    strings.TrimSpace(d.URL) != "",
		"migrations_dir": d.MigrationsDir,
		"schema":         d.Schema,
	}
}

func ValidateDatabase(cfg DatabaseConfig) error {
	if !cfg.Enabled {
		return nil
	}
	if strings.TrimSpace(cfg.Driver) == "" {
		return fmt.Errorf("database.driver is required when database is enabled")
	}
	if strings.TrimSpace(cfg.URL) == "" {
		return fmt.Errorf("database.url is required when database is enabled")
	}
	if _, err := url.Parse(cfg.URL); err != nil {
		return fmt.Errorf("database.url is invalid: %w", err)
	}
	if strings.TrimSpace(cfg.MigrationsDir) == "" {
		return fmt.Errorf("database.migrations_dir is required when database is enabled")
	}
	if strings.TrimSpace(cfg.Schema) == "" {
		return fmt.Errorf("database.schema is required when database is enabled")
	}
	return nil
}
