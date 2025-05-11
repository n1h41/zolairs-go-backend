package repositories

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"n1h41/zolaris-backend-app/internal/models"
)

type UserRepository struct {
	db          *dynamodb.Client
	userTable   string
	entityTable string
}

func NewUserRepository(dbClient *dynamodb.Client) *UserRepository {
	return &UserRepository{
		db:          dbClient,
		userTable:   "user_table",
		entityTable: "entityTable",
	}
}

func (r *UserRepository) WithTables(userTable, entityTable string) *UserRepository {
	r.userTable = userTable
	r.entityTable = entityTable
	return r
}

func (r *UserRepository) HasParentID(ctx context.Context, userID string) (bool, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(r.entityTable),
		IndexName:              aws.String("idType-index"),
		KeyConditionExpression: aws.String("idType = :idType"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":idType": &types.AttributeValueMemberS{Value: userID},
		},
	}

	result, err := r.db.Query(ctx, input)
	if err != nil {
		return false, err
	}

	if len(result.Items) == 0 {
		return false, nil
	}

	for _, item := range result.Items {
		_, exists := item["parentId"]
		if exists {
			return true, nil
		}
	}

	return false, nil
}

// UpdateUserDetails adds or updates the user details for a specific user in DynamoDB
func (r *UserRepository) UpdateUserDetails(ctx context.Context, userID string, details *models.UserDetails) error {
	// First, check if the user exists
	checkInput := &dynamodb.GetItemInput{
		TableName: aws.String(r.userTable),
		Key: map[string]types.AttributeValue{
			"user_id": &types.AttributeValueMemberS{Value: userID},
		},
		ProjectionExpression: aws.String("user_id"),
	}

	checkResult, err := r.db.GetItem(ctx, checkInput)
	if err != nil {
		return fmt.Errorf("error checking user existence: %w", err)
	}

	if len(checkResult.Item) == 0 {
		return fmt.Errorf("user with ID %s does not exist", userID)
	}

	// Create attribute values for the update expression
	expressionAttrValues := map[string]types.AttributeValue{
		":ud": &types.AttributeValueMemberM{
			Value: map[string]types.AttributeValue{
				"city":      &types.AttributeValueMemberS{Value: details.City},
				"country":   &types.AttributeValueMemberS{Value: details.Country},
				"email":     &types.AttributeValueMemberS{Value: details.Email},
				"firstName": &types.AttributeValueMemberS{Value: details.FirstName},
				"lastName":  &types.AttributeValueMemberS{Value: details.LastName},
				"phone":     &types.AttributeValueMemberS{Value: details.Phone},
				"region":    &types.AttributeValueMemberS{Value: details.Region},
				"street1":   &types.AttributeValueMemberS{Value: details.Street1},
				"zip":       &types.AttributeValueMemberS{Value: details.Zip},
			},
		},
	}

	// Add street2 only if it's not empty
	if details.Street2 != "" {
		expressionAttrValues[":ud"].(*types.AttributeValueMemberM).Value["street2"] = &types.AttributeValueMemberS{Value: details.Street2}
	}

	// Create update item input
	input := &dynamodb.UpdateItemInput{
		TableName: aws.String(r.userTable),
		Key: map[string]types.AttributeValue{
			"user_id": &types.AttributeValueMemberS{Value: userID},
		},
		UpdateExpression:          aws.String("SET userDetails = :ud"),
		ExpressionAttributeValues: expressionAttrValues,
	}

	// Execute the update
	_, err = r.db.UpdateItem(ctx, input)
	return err
}

// GetUserDetails retrieves the user details for a specific user from DynamoDB
func (r *UserRepository) GetUserDetails(ctx context.Context, userID string) (*models.UserDetails, error) {
	// Create get item input
	input := &dynamodb.GetItemInput{
		TableName: aws.String(r.userTable),
		Key: map[string]types.AttributeValue{
			"user_id": &types.AttributeValueMemberS{Value: userID},
		},
		ProjectionExpression: aws.String("userDetails"),
	}

	// Execute the get item operation
	result, err := r.db.GetItem(ctx, input)
	if err != nil {
		return nil, err
	}

	// Check if the item exists and has userDetails
	if len(result.Item) == 0 {
		return nil, nil // User not found
	}

	userDetailsAttr, exists := result.Item["userDetails"]
	if !exists {
		return nil, nil // User exists but has no details
	}

	// Parse the user details from the DynamoDB attribute
	userDetailsMap, ok := userDetailsAttr.(*types.AttributeValueMemberM)
	if !ok {
		return nil, fmt.Errorf("unexpected attribute type for userDetails")
	}

	userDetails := &models.UserDetails{}

	// Extract each field from the map, safely handling missing attributes
	if cityAttr, exists := userDetailsMap.Value["city"]; exists {
		if v, ok := cityAttr.(*types.AttributeValueMemberS); ok {
			userDetails.City = v.Value
		}
	}

	if countryAttr, exists := userDetailsMap.Value["country"]; exists {
		if v, ok := countryAttr.(*types.AttributeValueMemberS); ok {
			userDetails.Country = v.Value
		}
	}

	if emailAttr, exists := userDetailsMap.Value["email"]; exists {
		if v, ok := emailAttr.(*types.AttributeValueMemberS); ok {
			userDetails.Email = v.Value
		}
	}

	if firstNameAttr, exists := userDetailsMap.Value["firstName"]; exists {
		if v, ok := firstNameAttr.(*types.AttributeValueMemberS); ok {
			userDetails.FirstName = v.Value
		}
	}

	if lastNameAttr, exists := userDetailsMap.Value["lastName"]; exists {
		if v, ok := lastNameAttr.(*types.AttributeValueMemberS); ok {
			userDetails.LastName = v.Value
		}
	}

	if phoneAttr, exists := userDetailsMap.Value["phone"]; exists {
		if v, ok := phoneAttr.(*types.AttributeValueMemberS); ok {
			userDetails.Phone = v.Value
		}
	}

	if regionAttr, exists := userDetailsMap.Value["region"]; exists {
		if v, ok := regionAttr.(*types.AttributeValueMemberS); ok {
			userDetails.Region = v.Value
		}
	}

	if street1Attr, exists := userDetailsMap.Value["street1"]; exists {
		if v, ok := street1Attr.(*types.AttributeValueMemberS); ok {
			userDetails.Street1 = v.Value
		}
	}

	if street2Attr, exists := userDetailsMap.Value["street2"]; exists {
		if v, ok := street2Attr.(*types.AttributeValueMemberS); ok {
			userDetails.Street2 = v.Value
		}
	}

	if zipAttr, exists := userDetailsMap.Value["zip"]; exists {
		if v, ok := zipAttr.(*types.AttributeValueMemberS); ok {
			userDetails.Zip = v.Value
		}
	}

	return userDetails, nil
}
