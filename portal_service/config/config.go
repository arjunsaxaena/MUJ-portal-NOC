package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Port         string `mapstructure:"port"`
	Database     string `mapstructure:"database"`
	JwtSecretKey string `mapstructure:"jwt_secret_key"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("./portal_service/config")
	viper.AddConfigPath(".")

	fmt.Println("Current working directory:", viper.ConfigFileUsed())

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file: %s", err)
		return nil, err
	}

	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		log.Fatalf("Unable to decode into struct: %v", err)
		return nil, err
	}

	if config.JwtSecretKey == "" {
		return nil, fmt.Errorf("JWT secret key is required")
	}

	return &config, nil
}

func (c *Config) GetDatabaseURL() string {
	return c.Database
}
