package config

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/viper"
)

var cfg Config

func init() {
	err := LoadConfig([]string{"."}, &cfg)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
}

// Config represents the application configuration.
type Config struct {
	Service       ServiceConfig
	Server        ServerConfig
	Database      DatabaseConfig
	Cache         CacheConfig
	Env           string
	LogType       string
	ORM           ORMConfig
	Oidc          OidcConfig
	Messaging     MessagingConfig
	Logger        LoggerConfig
	Observability ObservabilityConfig
	Modules       []ModuleConfig
}

// ServiceConfig represents the service configuration.
type ServiceConfig struct {
	Project string
	Name    string
	Version string
}

// ServerConfig represents the server configuration.
type ServerConfig struct {
	Port int
}

// ORMConfig represents the ORM configuration.
type ORMConfig struct {
	MigrateDB bool
}

// DatabaseConfig represents the database configuration.
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
	Enable   bool
	Type     string
}

// CacheConfig represents the cache configuration.
type CacheConfig struct {
	Addrs    string
	User     string
	Password string
	Type     string
	Enable   bool
}

// OidcConfig represents the OIDC configuration.
type OidcConfig struct {
	Issuer       string
	ClientId     string
	ClientSecret string
	Enable       bool
}

// MessagingConfig represents the messaging configuration.
type MessagingConfig struct {
	Url      string
	UserName string
	Password string
	Enable   bool
	Type     string
}

// ObservabilityConfig represents the observability configuration.
type ObservabilityConfig struct {
	Tracing TracingConfig
	Metrics MetricsConfig
}

// TracingConfig represents the tracing configuration.
type TracingConfig struct {
	Enable   bool
	Type     string
	Endpoint string
}

// MetricsConfig represents the metrics configuration.
type MetricsConfig struct {
	Enable   bool
	Type     string
	Endpoint string
}

// LoggerConfig represents the logger configuration.
type LoggerConfig struct {
	Enable   bool
	Type     string
	Endpoint string
}

type ModuleConfig struct {
	Id     string
	Name   string
	Enable bool
}

// LoadConfig loads the configuration from the specified paths.
func LoadConfig(configPaths []string, pluginCF interface{}) error {
	viper.SetConfigName("application")
	viper.SetConfigType("yaml")

	for _, path := range configPaths {
		viper.AddConfigPath(path)
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("No config file found. Using default settings and environment variables.")
		} else {
			return fmt.Errorf("error reading config file: %v", err)
		}
	}

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	if err := viper.Unmarshal(pluginCF); err != nil {
		return fmt.Errorf("error unmarshal config: %v", err)
	}
	return nil
}

// GetConfig returns the application configuration.
func GetConfig() *Config {
	return &cfg
}
