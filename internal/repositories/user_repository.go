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
	input := &dynamodb.GetItemInput{
		TableName: aws.String(r.userTable),
		Key: map[string]types.AttributeValue{
			"idType": &types.AttributeValueMemberS{Value: userID},
		},
	}

	result, err := r.db.GetItem(ctx, input)
	if err != nil {
		return false, err
	}

	_, exists := result.Item["parentId"]
	return exists, nil
}
