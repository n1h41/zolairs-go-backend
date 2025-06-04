package dto

import (
	"time"
)

// Response is a standardized API response envelope
type Response struct {
	Success bool   `json:"success"`
	Data    any    `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
	Error   string `json:"error,omitempty"`
}

// ErrorResponse represents an API error response
type ErrorResponse struct {
	Success bool   `json:"success"`
	Error   string `json:"error"`
	Code    string `json:"code,omitempty"`
}

// ValidationError represents field validation errors
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// UserResponse represents user data in API responses
type UserResponse struct {
	ID        string         `json:"id"`
	Email     string         `json:"email"`
	FirstName *string        `json:"firstName"`
	LastName  *string        `json:"lastName"`
	Phone     *string        `json:"phone"`
	Address   *AddressOutput `json:"address"`
	ParentID  string         `json:"parentId,omitempty"`
	CreatedAt time.Time      `json:"createdAt"`
}

// AddressOutput represents address data in responses
type AddressOutput struct {
	Street1 string `json:"street1"`
	Street2 string `json:"street2,omitempty"`
	City    string `json:"city"`
	Region  string `json:"region"`
	Country string `json:"country"`
	Zip     string `json:"zip"`
}

// DeviceResponse represents device data in API responses
type DeviceResponse struct {
	DeviceID    string    `json:"deviceId"`
	DeviceName  string    `json:"deviceName"`
	Category    string    `json:"category,omitempty"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
}

// SensorDataResponse represents sensor readings in API responses
type SensorDataResponse struct {
	Timestamp   int64  `json:"timestamp"`
	Amperage    string `json:"amperage"`
	Temperature string `json:"temperature"`
	Humidity    string `json:"humidity"`
}

// CategoryResponse represents category data in API responses
type CategoryResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

// PaginatedResponse wraps list responses with pagination metadata
type PaginatedResponse struct {
	Items      any   `json:"items"`
	TotalItems int64 `json:"totalItems"`
	Page       int   `json:"page"`
	PageSize   int   `json:"pageSize"`
	TotalPages int   `json:"totalPages"`
}

// EntityResponse represents an entity in API responses
type EntityResponse struct {
	ID           string         `json:"id"`
	Name         string         `json:"name"`
	UserID       string         `json:"userId,omitempty"`
	CategoryID   string         `json:"categoryId"`
	CategoryName string         `json:"categoryName,omitempty"`
	CategoryType string         `json:"categoryType,omitempty"`
	ParentID     string         `json:"parentId,omitempty"`
	Details      map[string]any `json:"details,omitempty"`
	Depth        int            `json:"depth"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
}

// EntityHierarchyResponse represents an entity with its children in API responses
type EntityHierarchyResponse struct {
	EntityResponse
	Children []EntityHierarchyResponse `json:"children,omitempty"`
}

// EntityChildrenResponse represents a list of entity children
type EntityChildrenResponse struct {
	ParentID string            `json:"parentId"`
	Children []*EntityResponse `json:"children"`
	Count    int               `json:"count"`
}
