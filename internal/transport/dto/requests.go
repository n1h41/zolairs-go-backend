package dto

import "time"

// UserDetailsRequest represents a request to create or update user details
type UserDetailsRequest struct {
	Email     string `json:"email" validate:"required,email"`
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	Phone     string `json:"phone" validate:"required"`
	Street1   string `json:"street1" validate:"required"`
	Street2   string `json:"street2"`
	City      string `json:"city" validate:"required"`
	Region    string `json:"region" validate:"required"`
	Country   string `json:"country" validate:"required"`
	Zip       string `json:"zip" validate:"required"`
	ParentID  string `json:"parentId,omitempty"`
}

// DeviceRequest represents a request to add a new device
type DeviceRequest struct {
	DeviceID    string `json:"deviceId" validate:"required,min=3,max=50"`
	DeviceName  string `json:"deviceName" validate:"required,min=1,max=100"`
	Description string `json:"description,omitempty"`
	Category    string `json:"category,omitempty"`
}

// CategoryRequest represents a request to add a new category
type CategoryRequest struct {
	Name string `json:"name" validate:"required,min=2,max=50"`
	Type string `json:"type" validate:"required,min=2,max=50"`
}

// PolicyAttachRequest represents a request to attach an IoT policy
type PolicyAttachRequest struct {
	IdentityID string `json:"identityId" validate:"required"`
}

// SensorDataRequest represents a request to get device sensor data
type SensorDataRequest struct {
	DeviceMacID string `json:"deviceMacId" validate:"required"`
	Timestamp   int64  `json:"timestamp" validate:"required"`
	DateMode    string `json:"dateMode" validate:"required,oneof=hourly daily weekly monthly yearly"`
}

// TimeRange defines start and end times for data filtering
type TimeRange struct {
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
}

// PaginationParams contains pagination parameters
type PaginationParams struct {
	Page     int `json:"page" form:"page" default:"0"`
	PageSize int `json:"pageSize" form:"pageSize" default:"20"`
}

