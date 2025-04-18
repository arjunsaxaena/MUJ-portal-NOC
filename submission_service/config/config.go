package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port     string
	Database string
}

func LoadConfig() (*Config, error) {
	config := &Config{
		Port:     os.Getenv("SUBMISSION_SERVICE_PORT"),
		Database: os.Getenv("DB_URL"),
	}

	if config.Port == "" {
		config.Port = "8001"
	}

	if config.Database == "" {
		return nil, fmt.Errorf("database URL is required")
	}

	return config, nil
}

func (c *Config) GetDatabaseURL() string {
	return c.Database
}
