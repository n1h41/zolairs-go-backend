package repositories

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
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
