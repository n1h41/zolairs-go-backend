package aws

import (
	"context"
	"fmt"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/iot"
)

type Clients struct {
	DynamoDB *dynamodb.Client
	Iot      *iot.Client
}

func InitAWSClients(ctx context.Context) (*Clients, error) {
	awsCfg, err := awsconfig.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	return &Clients{
		DynamoDB: dynamodb.NewFromConfig(awsCfg),
		Iot:      iot.NewFromConfig(awsCfg),
	}, nil
}

func (c *Clients) GetIoTClient() *iot.Client {
	return c.Iot
}
