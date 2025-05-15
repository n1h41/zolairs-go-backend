package handlers

import (
	"log"

	"github.com/gin-gonic/gin"

	"n1h41/zolaris-backend-app/internal/middleware"
	"n1h41/zolaris-backend-app/internal/models"
	"n1h41/zolaris-backend-app/internal/services"
	"n1h41/zolaris-backend-app/internal/transport/response"
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
// @Param X-User-ID header string true "User ID"
// @Param device body models.AddDeviceRequest true "Device information"
// @Success 201 {object} dto.Response "Device added successfully"
// @Failure 400 {object} dto.ErrorResponse "Validation error"
// @Failure 401 {object} dto.ErrorResponse "User not authenticated"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Security ApiKeyAuth
// @Router /device/add [post]
func (h *AddDeviceHandler) HandleGin(c *gin.Context) {
	// Parse request body
	var request models.AddDeviceRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Printf("Error decoding request: %v", err)
		response.BadRequest(c, "Invalid request format")
		return
	}

	// Validate request
	validationErrs := utils.Validate(request)
	if validationErrs != nil {
		log.Printf("Validation errors: %s", utils.ValidationErrorsToString(validationErrs))
		response.ValidationErrors(c, utils.CreateDtoValidationErrors(validationErrs))
		return
	}

	// Get user ID from context (set by auth middleware)
	userID := middleware.GetUserIDFromGin(c)
	if userID == "" {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	// Call service to add device
	if err := h.deviceService.AddDevice(c.Request.Context(), request.DeviceID, request.DeviceName, userID); err != nil {
		log.Printf("Error adding device: %v", err)
		response.InternalError(c, "Failed to add device")
		return
	}

	response.Created(c, nil, "Device added successfully")
}
