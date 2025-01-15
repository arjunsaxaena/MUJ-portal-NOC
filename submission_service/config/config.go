package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Port     string `mapstructure:"port"`
	Database string `mapstructure:"database"`
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath("./config")

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

	return &config, nil
}

func (c *Config) GetDatabaseURL() string {
	return c.Database
}
