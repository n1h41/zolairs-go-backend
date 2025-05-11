package models

// UserDetailsRequest represents a request to update user details
type UserDetailsRequest struct {
	City      string `json:"city" validate:"required"`
	Country   string `json:"country" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	Phone     string `json:"phone" validate:"required"`
	Region    string `json:"region" validate:"required"`
	Street1   string `json:"street1" validate:"required"`
	Street2   string `json:"street2"`
	Zip       string `json:"zip" validate:"required"`
}

// AttachIotPolicyRequest represents a request to attach an IoT policy
type AttachIotPolicyRequest struct {
	IdentityId string `json:"identityId" validate:"required"`
}

// GetDeviceSensorDataRequest represents a request to get device sensor data
type GetDeviceSensorDataRequest struct {
	DeviceMacId string `json:"deviceMacId" validate:"required"`
	Timestamp   string `json:"timestamp" validate:"required"`
	DateMode    string `json:"dateMode" validate:"required,oneof=hourly daily weekly monthly yearly"`
}

// AddDeviceRequest represents a request to add a new device
type AddDeviceRequest struct {
	DeviceID   string `json:"deviceId" validate:"required,min=3,max=50"`
	DeviceName string `json:"deviceName" validate:"required,min=1,max=100"`
}

// AddCategoryRequest represents a request to add a new category
type AddCategoryRequest struct {
	Name string `json:"name" validate:"required,min=2,max=50"`
	Type string `json:"type" validate:"required,min=2,max=50"`
}



