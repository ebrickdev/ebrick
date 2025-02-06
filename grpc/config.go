package grpc

import (
	"fmt"

	"github.com/ebrickdev/ebrick/config"
)

type Config struct {
	Grpc GrpcServerConfig `yaml:"grpc"`
}

type GrpcServerConfig struct {
	Enabled bool   `yaml:"enabled"`
	Address string `yaml:"address"`
}

// GetConfig loads and returns the configuration
func GetConfig() (*Config, error) {
	var cfg Config
	if err := config.LoadConfig("application", []string{"."}, &cfg); err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}
	// Set default address if enabled and address is empty
	if cfg.Grpc.Enabled && cfg.Grpc.Address == "" {
		cfg.Grpc.Address = ":50051"
	}
	return &cfg, nil
}
