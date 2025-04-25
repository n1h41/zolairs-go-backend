package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config represents the application configuration
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	AWS      AWSConfig
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port int
}

// DatabaseConfig holds database-related configuration
type DatabaseConfig struct {
	DeviceTableName string
	DataTableName   string
	UserTableName   string
}

// AWSConfig holds AWS-related configuration
type AWSConfig struct {
	Region    string
	IoTPolicy string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*Config, error) {
	config := &Config{}

	// Server config
	portStr := getEnv("PORT", "8080")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("invalid PORT value: %v", err)
	}
	config.Server.Port = port

	// Database config
	config.Database.DeviceTableName = getEnv("DEVICE_TABLE_NAME", "devices")
	config.Database.DataTableName = getEnv("DATA_TABLE_NAME", "machine_data")
	config.Database.UserTableName = getEnv("USER_TABLE_NAME", "users")

	// AWS config
	config.AWS.Region = getEnv("AWS_REGION", "us-east-1")
	config.AWS.IoTPolicy = getEnv("IOT_POLICY_NAME", "DefaultIoTPolicy")

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

