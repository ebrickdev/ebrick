package config

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/viper"
)

var cfg Config

func init() {
	err := LoadConfig("application", []string{"."}, &cfg)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
}

// Config represents the application configuration.
type Config struct {
	Env    string
	Server ServerConfig
	Grpc   GrpcConfig
}
type ServerConfig struct {
	Port string
}

type GrpcConfig struct {
	Enable bool
}

// LoadConfig loads the configuration from the specified paths.
func LoadConfig(configName string, configPaths []string, data any) error {
	viper.SetConfigName(configName)
	viper.SetConfigType("yaml")

	for _, path := range configPaths {
		viper.AddConfigPath(path)
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("No config file found. Using default settings and environment variables.")
		} else {
			return fmt.Errorf("error reading config file %s: %v", configName, err)
		}
	}

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	if err := viper.Unmarshal(data); err != nil {
		return fmt.Errorf("error unmarshal config: %v", err)
	}
	return nil
}

// GetConfig returns the application configuration.
func GetAppConfig() *Config {
	return &cfg
}
