package db

import (
	"context"
	"log"
	"n1h41/zolaris-backend-app/internal/config"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/iot"
)

// Database holds the database clients and configuration
type Database struct {
	dynamoClient *dynamodb.Client
	iotClient    *iot.Client
	config       *config.Config
}

// NewDatabase creates and initializes database clients
func NewDatabase(cfg *config.Config) (*Database, error) {
	// Load AWS configuration
	log.Println("Initializing AWS clients...")
	awsCfg, err := awsconfig.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	// Create DynamoDB client
	dynamoClient := dynamodb.NewFromConfig(awsCfg)

	// Create IoT client
	iotClient := iot.NewFromConfig(awsCfg)

	return &Database{
		dynamoClient: dynamoClient,
		iotClient:    iotClient,
		config:       cfg,
	}, nil
}

// GetDynamoClient returns the DynamoDB client
func (db *Database) GetDynamoClient() *dynamodb.Client {
	return db.dynamoClient
}

// GetIoTClient returns the IoT client
func (db *Database) GetIoTClient() *iot.Client {
	return db.iotClient
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

