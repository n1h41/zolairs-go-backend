package repositories

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"n1h41/zolaris-backend-app/internal/domain"
)

// CategoryDBModel represents how the category is stored in the database
type CategoryDBModel struct {
	ID   string `dynamodbav:"id"`
	Name string `dynamodbav:"name"`
	Type string `dynamodbav:"type"`
}

// CategoryRepository handles all category-related database operations
type CategoryRepository struct {
	db            *dynamodb.Client
	categoryTable string
}

// NewCategoryRepository creates a new category repository instance
func NewCategoryRepository(dbClient *dynamodb.Client) *CategoryRepository {
	return &CategoryRepository{
		db:            dbClient,
		categoryTable: "categoryTable",
	}
}

// WithTable sets the table name for the repository
func (r *CategoryRepository) WithTable(categoryTable string) *CategoryRepository {
	r.categoryTable = categoryTable
	return r
}

// AddCategory adds a new category to the database
func (r *CategoryRepository) AddCategory(ctx context.Context, name, categoryType string) error {
	// Create item
	input := &dynamodb.PutItemInput{
		TableName: aws.String(r.categoryTable),
		Item: map[string]types.AttributeValue{
			"name": &types.AttributeValueMemberS{Value: name},
			"type": &types.AttributeValueMemberS{Value: categoryType},
		},
	}

	// Put item in DynamoDB
	_, err := r.db.PutItem(ctx, input)
	return err
}

// GetCategoryByName retrieves a category by its name
func (r *CategoryRepository) GetCategoryByName(ctx context.Context, name string) (*domain.Category, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(r.categoryTable),
		Key: map[string]types.AttributeValue{
			"name": &types.AttributeValueMemberS{Value: name},
		},
	}

	// Execute the query
	result, err := r.db.GetItem(ctx, input)
	if err != nil {
		return nil, err
	}

	// If item not found
	if result.Item == nil {
		return nil, nil
	}

	// Unmarshal the results
	var dbCategory CategoryDBModel
	err = attributevalue.UnmarshalMap(result.Item, &dbCategory)
	if err != nil {
		return nil, err
	}

	// Convert to domain model
	domainCategory := &domain.Category{
		ID:   dbCategory.ID,
		Name: dbCategory.Name,
		Type: dbCategory.Type,
	}

	return domainCategory, nil
}

// GetCategoriesByType retrieves all categories of a specific type
func (r *CategoryRepository) GetCategoriesByType(ctx context.Context, categoryType string) ([]*domain.Category, error) {
	input := &dynamodb.QueryInput{
		TableName:              aws.String(r.categoryTable),
		IndexName:              aws.String("TypeIndex"),
		KeyConditionExpression: aws.String("#type = :typeValue"),
		ExpressionAttributeNames: map[string]string{
			"#type": "type",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":typeValue": &types.AttributeValueMemberS{Value: categoryType},
		},
	}

	// Execute the query
	result, err := r.db.Query(ctx, input)
	if err != nil {
		return nil, err
	}

	// Unmarshal the results
	var dbCategories []CategoryDBModel
	err = attributevalue.UnmarshalListOfMaps(result.Items, &dbCategories)
	if err != nil {
		return nil, err
	}

	// Convert to domain model
	domainCategories := make([]*domain.Category, len(dbCategories))
	for i, cat := range dbCategories {
		domainCategories[i] = &domain.Category{
			ID:   cat.ID,
			Name: cat.Name,
			Type: cat.Type,
		}
	}

	return domainCategories, nil
}

// ListAllCategories retrieves all categories from the database
func (r *CategoryRepository) ListAllCategories(ctx context.Context) ([]*domain.Category, error) {
	// Create scan input
	input := &dynamodb.ScanInput{
		TableName: aws.String(r.categoryTable),
	}

	// Execute the scan
	result, err := r.db.Scan(ctx, input)
	if err != nil {
		return nil, err
	}

	// Unmarshal the results
	var dbCategories []CategoryDBModel
	err = attributevalue.UnmarshalListOfMaps(result.Items, &dbCategories)
	if err != nil {
		return nil, err
	}

	// Convert to domain model
	domainCategories := make([]*domain.Category, len(dbCategories))
	for i, cat := range dbCategories {
		domainCategories[i] = &domain.Category{
			ID:   cat.ID,
			Name: cat.Name,
			Type: cat.Type,
		}
	}

	return domainCategories, nil
}
