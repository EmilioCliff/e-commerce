package utils

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	DB_DRIVER string `mapstructure:"DB_DRIVER"`
}

// Reads config files and loads it into the Config Struct
func ReadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	err = viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return Config{}, fmt.Errorf("Config file not found: %w", err)
		} else {
			return Config{}, fmt.Errorf("Config file was found but another error was produced: %w", err)
		}
	}

	viper.Unmarshal(&config)
	return config, nil
}
