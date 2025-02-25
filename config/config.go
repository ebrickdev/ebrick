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
}
type ServerConfig struct {
	Port string
}

// LoadConfig loads the configuration from the specified paths.
func LoadConfig(configName string, configPaths []string, data any, defaults ...map[string]any) error {
	viper.SetConfigName(configName)
	viper.SetConfigType("yaml")

	// Apply default values if provided
	if len(defaults) > 0 {
		for key, value := range defaults[0] {
			viper.SetDefault(key, value)
		}
	}

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

func LoadConfigByKey(configName, key string, configPaths []string, data any, defaults ...map[string]any) error {
	// Set configuration file details
	viper.SetConfigName(configName)
	viper.SetConfigType("yaml")

	// Apply default values if provided
	if len(defaults) > 0 {
		for key, value := range defaults[0] {
			viper.SetDefault(key, value)
		}
	}

	// Set search paths for config files
	for _, path := range configPaths {
		viper.AddConfigPath(path)
	}

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.Sub(key).Unmarshal(data); err != nil {
		return fmt.Errorf("error unmarshal config: %v", err)
	}
	return nil
}
