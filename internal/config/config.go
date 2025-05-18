package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config represents the application configuration
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	AWS      AWSConfig
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port        int
	Environment string
	ExternalURL string
}

// DatabaseConfig holds database-related configuration
type DatabaseConfig struct {
	DeviceTableName  string
	DataTableName    string
	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDBName   string
	PostgresSSLMode  string
}

// AWSConfig holds AWS-related configuration
type AWSConfig struct {
	Region    string
	IoTPolicy string
}

// LoadEnv loads environment variables from .env files
func LoadEnv() error {
	// Try to load environment-specific .env file first
	environment := os.Getenv("ENVIRONMENT")
	if environment == "" {
		environment = "development" // Default to development
	}

	// Try environment specific file (.env.development, .env.production, etc)
	err := godotenv.Load(".env."+environment, ".env")
	if err != nil {
		// We'll just log the error and continue, as the .env file might not exist in all environments
		fmt.Printf("Warning: Error loading environment files: %v\n", err)
	}
	return err
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	// Try to load environment variables from .env files
	// We ignore the error as it's not fatal if .env files don't exist
	_ = LoadEnv()

	config := &Config{}

	// Server config
	portStr := getEnv("PORT", "8080")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("invalid PORT value: %v", err)
	}
	config.Server.Port = port
	config.Server.Environment = getEnv("ENVIRONMENT", "development")
	config.Server.ExternalURL = getEnv("EXTERNAL_URL", "http://localhost:8080")

	// Database config
	config.Database.DeviceTableName = getEnv("DEVICE_TABLE_NAME", "DEVICE_TABLE_NAME")
	config.Database.DataTableName = getEnv("DATA_TABLE_NAME", "machine_data_table")
	config.Database.PostgresHost = getEnv("POSTGRES_HOST", "localhost")
	config.Database.PostgresPort = getEnv("POSTGRES_PORT", "5432")
	config.Database.PostgresUser = getEnv("POSTGRES_USER", "postgres")
	config.Database.PostgresPassword = getEnv("POSTGRES_PASSWORD", "postgres")
	config.Database.PostgresDBName = getEnv("POSTGRES_DB_NAME", "postgres")
	config.Database.PostgresSSLMode = getEnv("POSTGRES_SSL_MODE", "disable")

	// AWS config
	config.AWS.Region = getEnv("AWS_REGION", "us-east-1")
	config.AWS.IoTPolicy = getEnv("IOT_POLICY_NAME", "IOT_POLICY_NAME")

	return config, nil
}

// LoadConfigWithPath loads configuration using explicit .env file paths
func LoadConfigWithPath(envPaths ...string) (*Config, error) {
	// Load specified env files
	if len(envPaths) > 0 {
		err := godotenv.Load(envPaths...)
		if err != nil {
			fmt.Printf("Warning: Error loading specified environment files: %v\n", err)
		}
	}

	config := &Config{}

	// Server config
	portStr := getEnv("PORT", "8080")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("invalid PORT value: %v", err)
	}
	config.Server.Port = port
	config.Server.Environment = getEnv("ENVIRONMENT", "development")
	config.Server.ExternalURL = getEnv("EXTERNAL_URL", "http://localhost:8080")

	// Database config
	config.Database.DeviceTableName = getEnv("DEVICE_TABLE_NAME", "machine_table")
	config.Database.DataTableName = getEnv("DATA_TABLE_NAME", "machine_data_table")
	config.Database.PostgresHost = getEnv("POSTGRES_HOST", "localhost")
	config.Database.PostgresPort = getEnv("POSTGRES_PORT", "5432")
	config.Database.PostgresUser = getEnv("POSTGRES_USER", "postgres")
	config.Database.PostgresPassword = getEnv("POSTGRES_PASSWORD", "postgres")
	config.Database.PostgresDBName = getEnv("POSTGRES_DB_NAME", "postgres")
	config.Database.PostgresSSLMode = getEnv("POSTGRES_SSL_MODE", "disable")

	// AWS config
	config.AWS.Region = getEnv("AWS_REGION", "us-east-1")
	config.AWS.IoTPolicy = getEnv("IOT_POLICY_NAME", "iot_p")

	return config, nil
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
