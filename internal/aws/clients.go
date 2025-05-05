package aws

import (
	"context"
	"fmt"
	"log"
	"n1h41/zolaris-backend-app/internal/config"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/iot"
)

type Clients struct {
	DynamoDB *dynamodb.Client
	Iot      *iot.Client
}

func InitAWSClients(ctx context.Context, cfg *config.Config) (*Clients, error) {
	log.Printf("Initializing AWS clients with region: %s", cfg.AWS.Region)

	opts := []func(*awsconfig.LoadOptions) error{}
	if cfg.AWS.Region != "" {
		opts = append(opts, awsconfig.WithRegion(cfg.AWS.Region))
	}

	if cfg.AWS.Profile != "" {
		opts = append(opts, awsconfig.WithSharedConfigProfile(cfg.AWS.Profile))
	}

	awsCfg, err := awsconfig.LoadDefaultConfig(ctx, opts...)
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
