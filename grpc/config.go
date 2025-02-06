package grpc

import (
	"errors"
	"fmt"

	"github.com/ebrickdev/ebrick/config"
)

type Config struct {
	Grpc GrpcServerConfig `yaml:"grpc"`
}

type GrpcServerConfig struct {
	Enable  bool   `yaml:"enable"`
	Address string `yaml:"address"`
}

// GetConfig loads and returns the configuration
func GetConfig() (*Config, error) {
	var cfg Config
	err := config.LoadConfig("application", []string{"."}, &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}
	// Validate the configuration
	if cfg.Grpc.Address == "" {
		return nil, errors.New("grpc_server.address is required when gRPC is defined")
	}
	return &cfg, nil
}
