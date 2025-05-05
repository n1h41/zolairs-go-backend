package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Helper function to unset all relevant environment variables
func unsetEnvVars() {
	vars := []string{
		"PORT", "ENVIRONMENT", "DEVICE_TABLE_NAME", "DATA_TABLE_NAME", 
		"USER_TABLE_NAME", "AWS_REGION", "IOT_POLICY_NAME",
	}
	for _, v := range vars {
		os.Unsetenv(v)
	}
}

func TestLoadConfig(t *testing.T) {
	// Backup existing environment variables to restore later
	origEnv := map[string]string{}
	vars := []string{
		"PORT", "ENVIRONMENT", "DEVICE_TABLE_NAME", "DATA_TABLE_NAME", 
		"USER_TABLE_NAME", "AWS_REGION", "IOT_POLICY_NAME",
	}
	
	for _, v := range vars {
		origEnv[v] = os.Getenv(v)
	}

	// Clean up environment after test
	defer func() {
		for k, v := range origEnv {
			if v == "" {
				os.Unsetenv(k)
			} else {
				os.Setenv(k, v)
			}
		}
	}()

	t.Run("LoadsDefaultValues", func(t *testing.T) {
		// Make sure environment is clean
		unsetEnvVars()
		
		// No environment variables set, should use defaults
		config, err := LoadConfig()
		require.NoError(t, err)
		
		// Verify defaults
		assert.Equal(t, 8080, config.Server.Port)
		assert.Equal(t, "development", config.Server.Environment)
		assert.Equal(t, "machine_table", config.Database.DeviceTableName)
		assert.Equal(t, "machine_data_table", config.Database.DataTableName)
		assert.Equal(t, "user_table", config.Database.UserTableName)
		assert.Equal(t, "us-east-1", config.AWS.Region)
		assert.Equal(t, "iot_p", config.AWS.IoTPolicy)
	})

	t.Run("LoadsFromEnvironmentVariables", func(t *testing.T) {
		// Make sure environment is clean
		unsetEnvVars()
		
		// Set environment variables
		os.Setenv("PORT", "9090")
		os.Setenv("ENVIRONMENT", "production")
		os.Setenv("DEVICE_TABLE_NAME", "prod_device_table")
		os.Setenv("DATA_TABLE_NAME", "prod_data_table")
		os.Setenv("USER_TABLE_NAME", "prod_user_table")
		os.Setenv("AWS_REGION", "eu-west-1")
		os.Setenv("IOT_POLICY_NAME", "prod_iot_policy")

		config, err := LoadConfig()
		require.NoError(t, err)
		
		// Verify environment variables were used
		assert.Equal(t, 9090, config.Server.Port)
		assert.Equal(t, "production", config.Server.Environment)
		assert.Equal(t, "prod_device_table", config.Database.DeviceTableName)
		assert.Equal(t, "prod_data_table", config.Database.DataTableName)
		assert.Equal(t, "prod_user_table", config.Database.UserTableName)
		assert.Equal(t, "eu-west-1", config.AWS.Region)
		assert.Equal(t, "prod_iot_policy", config.AWS.IoTPolicy)
		
		// Clean up after this test
		unsetEnvVars()
	})

	t.Run("LoadsFromDotEnvFile", func(t *testing.T) {
		// Make sure environment is clean
		unsetEnvVars()
		
		// Create temporary .env file
		tempDir := t.TempDir()
		envFile := filepath.Join(tempDir, ".env")
		envContent := `PORT=7070
ENVIRONMENT=staging
DEVICE_TABLE_NAME=staging_device_table
DATA_TABLE_NAME=staging_data_table
USER_TABLE_NAME=staging_user_table
AWS_REGION=ap-northeast-1
IOT_POLICY_NAME=staging_iot_policy`
		
		err := os.WriteFile(envFile, []byte(envContent), 0644)
		require.NoError(t, err)
		
		// Use the explicit path method for testing
		config, err := LoadConfigWithPath(envFile)
		require.NoError(t, err)
		
		// Verify .env file values were used
		assert.Equal(t, 7070, config.Server.Port)
		assert.Equal(t, "staging", config.Server.Environment)
		assert.Equal(t, "staging_device_table", config.Database.DeviceTableName)
		assert.Equal(t, "staging_data_table", config.Database.DataTableName)
		assert.Equal(t, "staging_user_table", config.Database.UserTableName)
		assert.Equal(t, "ap-northeast-1", config.AWS.Region)
		assert.Equal(t, "staging_iot_policy", config.AWS.IoTPolicy)
	})

	t.Run("LoadsFromEnvironmentSpecificDotEnvFile", func(t *testing.T) {
		// Make sure environment is clean
		unsetEnvVars()
		
		// Create temporary .env files
		tempDir := t.TempDir()
		
		// Create regular .env file
		envFile := filepath.Join(tempDir, ".env")
		envContent := `PORT=7070
ENVIRONMENT=testing
DEVICE_TABLE_NAME=default_device_table
DATA_TABLE_NAME=default_data_table
USER_TABLE_NAME=default_user_table
AWS_REGION=ap-northeast-1
IOT_POLICY_NAME=default_iot_policy`
		
		err := os.WriteFile(envFile, []byte(envContent), 0644)
		require.NoError(t, err)
		
		// Create environment-specific .env file
		envTestingFile := filepath.Join(tempDir, ".env.testing")
		envTestingContent := `PORT=6060
DEVICE_TABLE_NAME=testing_device_table
DATA_TABLE_NAME=testing_data_table
USER_TABLE_NAME=testing_user_table
AWS_REGION=us-west-2
IOT_POLICY_NAME=testing_iot_policy`
		
		err = os.WriteFile(envTestingFile, []byte(envTestingContent), 0644)
		require.NoError(t, err)
		
		// Set ENVIRONMENT to testing to trigger loading .env.testing
		os.Setenv("ENVIRONMENT", "testing")
		
		// Use the explicit path method for better testability
		config, err := LoadConfigWithPath(envTestingFile, envFile)
		require.NoError(t, err)
		
		// Verify environment-specific .env file values were used
		// and they override the regular .env values
		assert.Equal(t, 6060, config.Server.Port)
		assert.Equal(t, "testing", config.Server.Environment) // From ENVIRONMENT var + file
		assert.Equal(t, "testing_device_table", config.Database.DeviceTableName)
		assert.Equal(t, "testing_data_table", config.Database.DataTableName)
		assert.Equal(t, "testing_user_table", config.Database.UserTableName)
		assert.Equal(t, "us-west-2", config.AWS.Region)
		assert.Equal(t, "testing_iot_policy", config.AWS.IoTPolicy)
		
		// Clean up after this test
		unsetEnvVars()
	})

	t.Run("HandlesInvalidPortValue", func(t *testing.T) {
		// Make sure environment is clean
		unsetEnvVars()
		
		os.Setenv("PORT", "invalid_port")
		
		_, err := LoadConfig()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid PORT value")
		
		// Clean up after this test
		unsetEnvVars()
	})
}

func TestGetEnv(t *testing.T) {
	// Save original value to restore after test
	origValue := os.Getenv("TEST_ENV_VAR")
	defer func() {
		if origValue == "" {
			os.Unsetenv("TEST_ENV_VAR")
		} else {
			os.Setenv("TEST_ENV_VAR", origValue)
		}
	}()

	// Test when environment variable is not set
	os.Unsetenv("TEST_ENV_VAR")
	result := getEnv("TEST_ENV_VAR", "default_value")
	assert.Equal(t, "default_value", result)

	// Test when environment variable is set
	os.Setenv("TEST_ENV_VAR", "custom_value")
	result = getEnv("TEST_ENV_VAR", "default_value")
	assert.Equal(t, "custom_value", result)

	// Test when environment variable is set but empty
	os.Setenv("TEST_ENV_VAR", "")
	result = getEnv("TEST_ENV_VAR", "default_value")
	assert.Equal(t, "default_value", result)
}
