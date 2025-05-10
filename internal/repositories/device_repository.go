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
	db           *dynamodb.Client
	deviceTable  string
	machineTable string
}

// NewDeviceRepository creates a new device repository instance
func NewDeviceRepository(dbClient *dynamodb.Client) *DeviceRepository {
	return &DeviceRepository{
		db:           dbClient,
		deviceTable:  "machine_table",
		machineTable: "machine_data_table",
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
			"mac_address": &types.AttributeValueMemberS{Value: deviceID},
			"user_id":     &types.AttributeValueMemberS{Value: userID},
			"device_name": &types.AttributeValueMemberS{Value: deviceName},
		},
	}

	// Put item in DynamoDB
	_, err := r.db.PutItem(ctx, input)
	return err
}

// GetDevicesByUserID retrieves all devices for a specific user
func (r *DeviceRepository) GetDevicesByUserID(ctx context.Context, userID string) ([]models.DeviceResponse, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(r.deviceTable),
		IndexName:              aws.String("user_id-index"),
		KeyConditionExpression: aws.String("user_id = :userId"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":userId": &types.AttributeValueMemberS{Value: userID},
		},
	}

	// Execute the query
	result, err := r.db.Query(ctx, input)
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
	log.Printf("Table name: %s", r.machineTable)

	input := &dynamodb.QueryInput{
		TableName:              aws.String(r.machineTable),
		KeyConditionExpression: aws.String("mac_id = :macId AND #ts BETWEEN :startTime AND :endTime"),
		ExpressionAttributeNames: map[string]string{
			"#ts": "timestamp",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":macId":     &types.AttributeValueMemberS{Value: macID},
			":startTime": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", startTime)},
			":endTime":   &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", endTime)},
		},
	}

	log.Printf("Querying sensor data for device %s from %d to %d", macID, startTime, endTime)

	result, err := r.db.Query(ctx, input)
	if err != nil {
		return nil, err
	}

	var sensorData []models.SensorData
	err = attributevalue.UnmarshalListOfMaps(result.Items, &sensorData)
	if err != nil {
		return nil, err
	}

	return sensorData, nil
}
