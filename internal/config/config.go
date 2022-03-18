package config

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	DB struct {
		Driver string `mapstructure:"driver"`
		Source string `mapstructure:"source"`
	}

	Servers struct {
		Grpc string `mapstructure:"grpc"`
		Http string `mapstructure:"http"`
	}
}

// const envPrefix = "ORDERS"

var ErrNoConfig = errors.New("no config set")

func LoadConfig() (*Config, error) {
	// viper.SetEnvPrefix(envPrefix)
	viper.AutomaticEnv()

	config := viper.GetString("config")
	if config == "" {
		return nil, ErrNoConfig
	}

	// viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.SetConfigName(config)

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("viper.ReadInConfig: %w", err)
	}

	var c Config

	if err := viper.Unmarshal(&c); err != nil {
		return nil, fmt.Errorf("viper.Unmarshal: %w", err)
	}

	return &c, nil
}
