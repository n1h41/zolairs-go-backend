package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"n1h41/zolaris-backend-app/internal/middleware"
	"n1h41/zolaris-backend-app/internal/models"
	"n1h41/zolaris-backend-app/internal/services"
	transport_gin "n1h41/zolaris-backend-app/internal/transport/gin"
	transport_http "n1h41/zolaris-backend-app/internal/transport/http"
	"n1h41/zolaris-backend-app/internal/utils"
)

// AddDeviceHandler handles requests to add a new device
type AddDeviceHandler struct {
	deviceService *services.DeviceService
}

// NewAddDeviceHandler creates a new AddDeviceHandler
func NewAddDeviceHandler(deviceService *services.DeviceService) *AddDeviceHandler {
	return &AddDeviceHandler{deviceService: deviceService}
}

// ServeHTTP implements http.Handler interface (for backward compatibility)
func (h *AddDeviceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Only allow POST method
	if r.Method != http.MethodPost {
		transport_http.SendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Parse request body
	var request models.AddDeviceRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Printf("Error decoding request: %v", err)
		transport_http.SendBadRequestError(w, "Invalid request format")
		return
	}

	// Validate request
	validationErrs := utils.Validate(request)
	if validationErrs != nil {
		log.Printf("Validation errors: %s", utils.ValidationErrorsToString(validationErrs))
		transport_http.SendBadRequestError(w, utils.CreateValidationError(validationErrs))
		return
	}

	// Get user ID from context (set by auth middleware)
	userID := middleware.GetUserID(r)
	if userID == "" {
		transport_http.SendError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	// Call service to add device
	if err := h.deviceService.AddDevice(r.Context(), request.DeviceID, request.DeviceName, userID); err != nil {
		log.Printf("Error adding device: %v", err)
		transport_http.SendError(w, http.StatusInternalServerError, "Failed to add device")
		return
	}

	transport_http.SendResponse(w, http.StatusCreated, "Device added successfully")
}

// HandleGin handles requests using Gin framework
func (h *AddDeviceHandler) HandleGin(c *gin.Context) {
	// Parse request body
	var request models.AddDeviceRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Printf("Error decoding request: %v", err)
		transport_gin.SendBadRequestError(c, "Invalid request format")
		return
	}

	// Validate request
	validationErrs := utils.Validate(request)
	if validationErrs != nil {
		log.Printf("Validation errors: %s", utils.ValidationErrorsToString(validationErrs))
		transport_gin.SendBadRequestError(c, utils.CreateValidationError(validationErrs))
		return
	}

	// Get user ID from context (set by auth middleware)
	userID := middleware.GetUserIDFromGin(c)
	if userID == "" {
		transport_gin.SendError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	// Call service to add device
	if err := h.deviceService.AddDevice(c.Request.Context(), request.DeviceID, request.DeviceName, userID); err != nil {
		log.Printf("Error adding device: %v", err)
		transport_gin.SendError(c, http.StatusInternalServerError, "Failed to add device")
		return
	}

	transport_gin.SendResponse(c, http.StatusCreated, "Device added successfully")
}

