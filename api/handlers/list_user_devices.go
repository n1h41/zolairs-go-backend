package handlers

import (
	"log"
	"n1h41/zolaris-backend-app/internal/middleware"
	"n1h41/zolaris-backend-app/internal/services"
	transport "n1h41/zolaris-backend-app/internal/transport/http"
	"net/http"
)

// ListUserDevicesHandler handles requests to list devices for a user
type ListUserDevicesHandler struct {
	deviceService *services.DeviceService
}

// NewListUserDevicesHandler creates a new ListUserDevicesHandler
func NewListUserDevicesHandler(deviceService *services.DeviceService) *ListUserDevicesHandler {
	return &ListUserDevicesHandler{deviceService: deviceService}
}

// ServeHTTP implements http.Handler interface
func (h *ListUserDevicesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Only allow GET method
	if r.Method != http.MethodGet {
		transport.SendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Get user ID from context (set by auth middleware)
	userID := middleware.GetUserID(r)
	if userID == "" {
		transport.SendError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	// Call service to get user devices
	devices, err := h.deviceService.GetUserDevices(r.Context(), userID)
	if err != nil {
		log.Printf("Error getting user devices: %v", err)
		transport.SendError(w, http.StatusInternalServerError, "Failed to retrieve user devices")
		return
	}

	transport.SendResponse(w, http.StatusOK, devices)
}

