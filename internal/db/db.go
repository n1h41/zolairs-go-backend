package db

import (
	"n1h41/zolaris-backend-app/internal/config"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// Database holds the database clients and configuration
type Database struct {
	dynamoClient *dynamodb.Client
	config       *config.Config
}

// NewDatabase creates and initializes database clients
func NewDatabase(db *dynamodb.Client, cfg *config.Config) (*Database, error) {
	return &Database{
		dynamoClient: db,
		config:       cfg,
	}, nil
}

func (db *Database) GetDynamoClient() *dynamodb.Client {
	return db.dynamoClient
}

// GetDeviceTableName returns the name of the device table
func (db *Database) GetDeviceTableName() string {
	return db.config.Database.DeviceTableName
}

// GetMachineDataTableName returns the name of the machine data table
func (db *Database) GetMachineDataTableName() string {
	return db.config.Database.DataTableName
}

// GetUserTableName returns the name of the user table
func (db *Database) GetUserTableName() string {
	return db.config.Database.UserTableName
}
