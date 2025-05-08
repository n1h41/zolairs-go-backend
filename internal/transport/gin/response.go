package gin

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response is the standard successful response structure
// @Description Standard API success response format
type Response struct {
	// Indicates success status
	Status bool `json:"status" example:"true"`
	// Optional success message
	Message string `json:"message,omitempty" example:"Operation successful"`
	// Optional data payload
	Data any `json:"data,omitempty"`
}

// ErrorResponse is the standard error response structure
// @Description Standard API error response format
type ErrorResponse struct {
	// Always false for errors
	Status bool `json:"status" example:"false"`
	// Error message
	Message string `json:"message" example:"Something went wrong"`
}

// SendResponse sends a successful JSON response
func SendResponse(c *gin.Context, statusCode int, data any) {
	resp := gin.H{
		"status": true,
	}

	// Handle different types of data
	switch v := data.(type) {
	case string:
		resp["message"] = v
	case any:
		resp["data"] = v
	}

	c.JSON(statusCode, resp)
}

// SendError sends an error JSON response
func SendError(c *gin.Context, statusCode int, err any) {
	log.Printf("Error: %v", err)

	errorResp := gin.H{
		"status": false,
	}

	// Handle different types of errors
	switch e := err.(type) {
	case error:
		errorResp["message"] = e.Error()
	case string:
		errorResp["message"] = e
	default:
		errorResp["message"] = "An unknown error occurred"
	}

	c.JSON(statusCode, errorResp)
	c.Abort()
}

// SendBadRequestError sends a 400 Bad Request error response
func SendBadRequestError(c *gin.Context, err any) {
	SendError(c, http.StatusBadRequest, err)
}

