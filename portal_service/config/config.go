package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port         string
	Database     string
	JwtSecretKey string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
		return nil, err
	}

	config := &Config{
		Port:         os.Getenv("PORTAL_SERVICE_PORT"),
		Database:     os.Getenv("DB_URL"),
		JwtSecretKey: os.Getenv("JWT_SECRET_KEY"),
	}

	if config.JwtSecretKey == "" {
		return nil, fmt.Errorf("JWT secret key is required")
	}

	return config, nil
}

func (c *Config) GetDatabaseURL() string {
	return c.Database
}
