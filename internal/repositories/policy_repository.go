package repositories

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/iot"
)

// PolicyRepository handles all IoT policy-related operations
type PolicyRepository struct {
	iotClient *iot.Client
}

// NewPolicyRepository creates a new policy repository instance
func NewPolicyRepository(iotClient *iot.Client) *PolicyRepository {
	return &PolicyRepository{iotClient: iotClient}
}

// AttachPolicy attaches an IoT policy to an identity
func (r *PolicyRepository) AttachPolicy(ctx context.Context, policyName, identityID string) error {
	// Create attach policy input
	input := &iot.AttachPolicyInput{
		PolicyName: &policyName,
		Target:     &identityID,
	}

	// Attach policy
	_, err := r.iotClient.AttachPolicy(ctx, input)
	return err
}
