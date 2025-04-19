package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port         string
	Database     string
	JwtSecretKey string
}

func LoadConfig() (*Config, error) {
	config := &Config{
		Port:         os.Getenv("PORTAL_SERVICE_PORT"),
		Database:     os.Getenv("DB_URL"),
		JwtSecretKey: os.Getenv("JWT_SECRET_KEY"),
	}

	if config.Port == "" {
		config.Port = "8002"
	}

	if config.Database == "" {
		return nil, fmt.Errorf("database URL is required")
	}

	return config, nil
}

func (c *Config) GetDatabaseURL() string {
	return c.Database
}
