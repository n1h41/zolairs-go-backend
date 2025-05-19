package repositories

import (
	"context"

	"n1h41/zolaris-backend-app/internal/domain"
)

// UserRepositoryInterface defines the operations for user data
type UserRepositoryInterface interface {
	GetUserIdByCognitoId(ctx context.Context, cId string) (string, error)
	GetUserByID(ctx context.Context, userID string) (*domain.User, error)
	CreateUser(ctx context.Context, user *domain.User) error
	UpdateUser(ctx context.Context, user *domain.User) error
	CheckHasParentID(ctx context.Context, userID string) (bool, error)
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	GetChildUsers(ctx context.Context, parentID string) ([]*domain.User, error)
}

// DeviceRepositoryInterface defines the operations for device data
type DeviceRepositoryInterface interface {
	AddDevice(ctx context.Context, device *domain.Device) error
	GetDevicesByUserID(ctx context.Context, userID string) ([]*domain.Device, error)
	GetSensorData(ctx context.Context, macID string, startTime, endTime int64) ([]*domain.SensorReading, error)
}

// CategoryRepositoryInterface defines the operations for category data
type CategoryRepositoryInterface interface {
	AddCategory(ctx context.Context, category *domain.Category) error
	GetCategoriesByType(ctx context.Context, categoryType string) ([]*domain.Category, error)
	ListAllCategories(ctx context.Context) ([]*domain.Category, error)
}

// PolicyRepositoryInterface defines the operations for policy data
type PolicyRepositoryInterface interface {
	AttachPolicy(ctx context.Context, identityId, policyName string) error
}
