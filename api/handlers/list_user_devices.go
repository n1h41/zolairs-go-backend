package handlers

import (
	"log"

	"github.com/gin-gonic/gin"

	"n1h41/zolaris-backend-app/internal/middleware"
	"n1h41/zolaris-backend-app/internal/services"
	"n1h41/zolaris-backend-app/internal/transport/response"
)

// ListUserDevicesHandler handles requests to list devices for a user
type ListUserDevicesHandler struct {
	deviceService *services.DeviceService
}

// NewListUserDevicesHandler creates a new ListUserDevicesHandler
func NewListUserDevicesHandler(deviceService *services.DeviceService) *ListUserDevicesHandler {
	return &ListUserDevicesHandler{deviceService: deviceService}
}

// HandleGin handles requests using Gin framework
// @Summary List user devices
// @Description Get all devices registered to the authenticated user
// @Tags Device Management
// @Accept json
// @Produce json
// @Param X-User-ID header string true "User ID"
// @Success 200 {array} models.DeviceResponse "List of user devices"
// @Failure 401 {object} dto.ErrorResponse "User not authenticated"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Security ApiKeyAuth
// @Router /user/devices [get]
func (h *ListUserDevicesHandler) HandleGin(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID := middleware.GetUserIDFromGin(c)
	if userID == "" {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	// Call service to get user devices
	devices, err := h.deviceService.GetUserDevices(c.Request.Context(), userID)
	if err != nil {
		log.Printf("Error getting user devices: %v", err)
		response.InternalError(c, "Failed to retrieve user devices")
		return
	}

	response.OK(c, devices, "Devices retrieved successfully")
}
