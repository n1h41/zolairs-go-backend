package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"n1h41/zolaris-backend-app/internal/middleware"
	"n1h41/zolaris-backend-app/internal/services"
	transport_gin "n1h41/zolaris-backend-app/internal/transport/gin"
	transport_http "n1h41/zolaris-backend-app/internal/transport/http"
)

// ListUserDevicesHandler handles requests to list devices for a user
type ListUserDevicesHandler struct {
	deviceService *services.DeviceService
}

// NewListUserDevicesHandler creates a new ListUserDevicesHandler
func NewListUserDevicesHandler(deviceService *services.DeviceService) *ListUserDevicesHandler {
	return &ListUserDevicesHandler{deviceService: deviceService}
}

// ServeHTTP implements http.Handler interface (for backward compatibility)
func (h *ListUserDevicesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Only allow GET method
	if r.Method != http.MethodGet {
		transport_http.SendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Get user ID from context (set by auth middleware)
	userID := middleware.GetUserID(r)
	if userID == "" {
		transport_http.SendError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	// Call service to get user devices
	devices, err := h.deviceService.GetUserDevices(r.Context(), userID)
	if err != nil {
		log.Printf("Error getting user devices: %v", err)
		transport_http.SendError(w, http.StatusInternalServerError, "Failed to retrieve user devices")
		return
	}

	transport_http.SendResponse(w, http.StatusOK, devices)
}

// HandleGin handles requests using Gin framework
func (h *ListUserDevicesHandler) HandleGin(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID := middleware.GetUserIDFromGin(c)
	if userID == "" {
		transport_gin.SendError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	// Call service to get user devices
	devices, err := h.deviceService.GetUserDevices(c.Request.Context(), userID)
	if err != nil {
		log.Printf("Error getting user devices: %v", err)
		transport_gin.SendError(c, http.StatusInternalServerError, "Failed to retrieve user devices")
		return
	}

	transport_gin.SendResponse(c, http.StatusOK, devices)
}
