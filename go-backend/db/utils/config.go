package utils

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DBDriver             string        `mapstructure:"DB_DRIVER"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	TokenFooter          string        `mapstructure:"TOKEN_FOOTER"`
}

// Reads config files and loads it into the Config Struct
func ReadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	err = viper.ReadInConfig()
	// errors.Is(err, viper.ConfigFileNotFoundError{})
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return Config{}, fmt.Errorf("Config file not found: %w", err)
		}
		return Config{}, fmt.Errorf("Config file was found but another error was produced: %w", err)
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return Config{}, err
	}
	return config, nil
}
