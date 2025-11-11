// Package config provides configuration management for the OSE Advisory Service.
//
// Configuration is loaded from multiple sources in order of precedence:
//   1. Command-line flags (highest precedence)
//   2. Environment variables (OSE_* prefix)
//   3. Config file (/etc/ose/config.yaml or specified via flag)
//   4. Default values (lowest precedence)
//
// This layered approach enables:
//   - Local development (defaults work out of the box)
//   - Container deployment (env vars from K8s ConfigMap)
//   - Production flexibility (config files for complex settings)
package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds all configuration for the Advisory Service.
type Config struct {
	// Server configuration
	GRPCAddress string // Address for gRPC server (e.g., ":50051")
	HTTPAddress string // Address for HTTP server (health checks, metrics)

	// Backend services (Week 5+)
	Neo4jURI      string // Neo4j connection string
	Neo4jUser     string // Neo4j username
	Neo4jPassword string // Neo4j password

	// Template engine (Week 6-8)
	TemplatePath string // Path to template directory

	// Observability
	LogLevel      string        // Log level: debug, info, warn, error
	EnableMetrics bool          // Enable Prometheus metrics endpoint
	EnableTracing bool          // Enable OpenTelemetry tracing
	TracingEndpoint string      // OTLP exporter endpoint

	// Performance tuning
	MaxConcurrentRequests int           // Max concurrent gRPC requests
	RequestTimeout        time.Duration // Timeout for individual requests
	ShutdownTimeout       time.Duration // Graceful shutdown timeout
}

// Load reads configuration from environment and returns a Config struct.
//
// Environment Variables:
//   OSE_GRPC_ADDRESS       - gRPC server address (default: ":50051")
//   OSE_HTTP_ADDRESS       - HTTP server address (default: ":8080")
//   OSE_NEO4J_URI          - Neo4j connection URI
//   OSE_NEO4J_USER         - Neo4j username
//   OSE_NEO4J_PASSWORD     - Neo4j password
//   OSE_TEMPLATE_PATH      - Template directory path
//   OSE_LOG_LEVEL          - Logging level (default: "info")
//   OSE_MAX_CONCURRENT     - Max concurrent requests (default: 100)
//   OSE_REQUEST_TIMEOUT    - Request timeout in seconds (default: 30)
func Load() (*Config, error) {
	cfg := &Config{
		// Server defaults
		GRPCAddress: getEnv("OSE_GRPC_ADDRESS", ":50051"),
		HTTPAddress: getEnv("OSE_HTTP_ADDRESS", ":8080"),

		// Backend defaults (will be required in Week 5)
		Neo4jURI:      getEnv("OSE_NEO4J_URI", "bolt://localhost:7687"),
		Neo4jUser:     getEnv("OSE_NEO4J_USER", "neo4j"),
		Neo4jPassword: getEnv("OSE_NEO4J_PASSWORD", ""),

		// Template defaults
		TemplatePath: getEnv("OSE_TEMPLATE_PATH", "./templates"),

		// Observability defaults
		LogLevel:        getEnv("OSE_LOG_LEVEL", "info"),
		EnableMetrics:   getEnvBool("OSE_ENABLE_METRICS", true),
		EnableTracing:   getEnvBool("OSE_ENABLE_TRACING", false),
		TracingEndpoint: getEnv("OSE_TRACING_ENDPOINT", "localhost:4317"),

		// Performance defaults
		MaxConcurrentRequests: getEnvInt("OSE_MAX_CONCURRENT", 100),
		RequestTimeout:        time.Duration(getEnvInt("OSE_REQUEST_TIMEOUT", 30)) * time.Second,
		ShutdownTimeout:       30 * time.Second,
	}

	// Validate required fields
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return cfg, nil
}

// Validate checks that all required configuration is present and valid.
func (c *Config) Validate() error {
	// In Week 3, most validations are lenient (backends not yet required)
	// Week 5+ will require Neo4j credentials

	if c.GRPCAddress == "" {
		return fmt.Errorf("GRPC address cannot be empty")
	}

	if c.HTTPAddress == "" {
		return fmt.Errorf("HTTP address cannot be empty")
	}

	if c.MaxConcurrentRequests <= 0 {
		return fmt.Errorf("max concurrent requests must be positive, got %d", c.MaxConcurrentRequests)
	}

	if c.RequestTimeout <= 0 {
		return fmt.Errorf("request timeout must be positive, got %v", c.RequestTimeout)
	}

	// Validate log level
	validLogLevels := map[string]bool{
		"debug": true, "info": true, "warn": true, "error": true,
	}
	if !validLogLevels[c.LogLevel] {
		return fmt.Errorf("invalid log level: %s (must be debug, info, warn, or error)", c.LogLevel)
	}

	return nil
}

// getEnv retrieves an environment variable with a default fallback.
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvInt retrieves an environment variable as integer with default fallback.
func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getEnvBool retrieves an environment variable as boolean with default fallback.
func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}
