package repositories

import (
	"context"
	"fmt"
	"log"
	"n1h41/zolaris-backend-app/internal/models"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// DeviceRepository handles all device-related database operations
type DeviceRepository struct {
	db            *dynamodb.Client
	deviceTable   string
	machineTable  string
}

// NewDeviceRepository creates a new device repository instance
func NewDeviceRepository(dbClient *dynamodb.Client) *DeviceRepository {
	return &DeviceRepository{
		db:           dbClient,
		deviceTable:  "devices",
		machineTable: "machine_data",
	}
}

// WithTables sets the table names for the repository
func (r *DeviceRepository) WithTables(deviceTable, machineTable string) *DeviceRepository {
	r.deviceTable = deviceTable
	r.machineTable = machineTable
	return r
}

// AddDevice adds a new device to the database
func (r *DeviceRepository) AddDevice(ctx context.Context, deviceID, deviceName, userID string) error {
	// Create item
	input := &dynamodb.PutItemInput{
		TableName: aws.String(r.deviceTable),
		Item: map[string]types.AttributeValue{
			"deviceId":   &types.AttributeValueMemberS{Value: deviceID},
			"userId":     &types.AttributeValueMemberS{Value: userID},
			"deviceName": &types.AttributeValueMemberS{Value: deviceName},
		},
	}

	// Put item in DynamoDB
	_, err := r.db.PutItem(ctx, input)
	return err
}

// GetDevicesByUserID retrieves all devices for a specific user
func (r *DeviceRepository) GetDevicesByUserID(ctx context.Context, userID string) ([]models.DeviceResponse, error) {
	// Create query statement
	statement := fmt.Sprintf(`SELECT * FROM "%s" WHERE "userId" = ?`, r.deviceTable)

	// Execute the PartiQL query with parameters
	input := &dynamodb.ExecuteStatementInput{
		Statement: aws.String(statement),
		Parameters: []types.AttributeValue{
			&types.AttributeValueMemberS{Value: userID},
		},
	}

	// Execute the query
	result, err := r.db.ExecuteStatement(ctx, input)
	if err != nil {
		return nil, err
	}

	// Unmarshal the results
	var devices []models.DeviceResponse
	err = attributevalue.UnmarshalListOfMaps(result.Items, &devices)
	if err != nil {
		return nil, err
	}

	return devices, nil
}

// GetSensorData retrieves sensor data for a specific device within a time range
func (r *DeviceRepository) GetSensorData(ctx context.Context, macID string, startTime, endTime int64) ([]models.SensorData, error) {
	// Create PartiQL statement
	statement := fmt.Sprintf(`SELECT * FROM "%s" WHERE "mac_id" = '%s' AND "timestamp" >= %d AND "timestamp" <= %d`,
		r.machineTable, macID, startTime, endTime)

	log.Printf("Executing PartiQL statement: %s", statement)

	// Execute the PartiQL query
	input := &dynamodb.ExecuteStatementInput{
		Statement: aws.String(statement),
	}

	// Execute the PartiQL statement
	result, err := r.db.ExecuteStatement(ctx, input)
	if err != nil {
		return nil, err
	}

	// Unmarshal the results
	var sensorData []models.SensorData
	err = attributevalue.UnmarshalListOfMaps(result.Items, &sensorData)
	if err != nil {
		return nil, err
	}

	return sensorData, nil
}

