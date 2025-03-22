package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port     string
	Database string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
		return nil, err
	}

	config := &Config{
		Port:     os.Getenv("SUBMISSION_SERVICE_PORT"),
		Database: os.Getenv("DB_URL"),
	}

	if config.Port == "" {
		return nil, fmt.Errorf("submission service port is required")
	}

	return config, nil
}

func (c *Config) GetDatabaseURL() string {
	return c.Database
}
