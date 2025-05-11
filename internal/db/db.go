package db

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"n1h41/zolaris-backend-app/internal/config"
)

// Database holds the database clients and configuration
type Database struct {
	dynamoClient *dynamodb.Client
	postgres     *PostgresDB
	config       *config.Config
}

// NewDatabase creates and initializes database clients
func NewDatabase(ctx context.Context, db *dynamodb.Client, cfg *config.Config) (*Database, error) {
	/* postgres, err := NewPostgresDB(ctx, cfg)
	if err != nil {
		return nil, err
	} */

	return &Database{
		// postgres:     postgres,
		dynamoClient: db,
		config:       cfg,
	}, nil
}

// GetPostgresDB returns the PostgreSQL database instance
func (db *Database) GetPostgresDB() *PostgresDB {
	return db.postgres
}

// Close releases all database resources
func (db *Database) Close() {
	if db.postgres != nil {
		db.postgres.Close()
	}
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
