package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"n1h41/zolaris-backend-app/internal/middleware"
	"n1h41/zolaris-backend-app/internal/models"
	"n1h41/zolaris-backend-app/internal/services"
	transport_gin "n1h41/zolaris-backend-app/internal/transport/gin"
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

// HandleGin handles requests using Gin framework
// @Summary Add a new device
// @Description Register a new IoT device for the authenticated user
// @Tags Device Management
// @Accept json
// @Produce json
// @Param device body models.AddDeviceRequest true "Device information"
// @Success 201 {object} transport_gin.Response "Device added successfully"
// @Failure 400 {object} transport_gin.ErrorResponse "Validation error"
// @Failure 401 {object} transport_gin.ErrorResponse "User not authenticated"
// @Failure 500 {object} transport_gin.ErrorResponse "Internal server error"
// @Security ApiKeyAuth
// @Router /device/add [post]
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
