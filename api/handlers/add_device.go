package handlers

import (
	"encoding/json"
	"log"
	"n1h41/zolaris-backend-app/internal/middleware"
	"n1h41/zolaris-backend-app/internal/models"
	"n1h41/zolaris-backend-app/internal/services"
	transport "n1h41/zolaris-backend-app/internal/transport/http"
	"n1h41/zolaris-backend-app/internal/utils"
	"net/http"
)

// AddDeviceHandler handles requests to add a new device
type AddDeviceHandler struct {
	deviceService *services.DeviceService
}

// NewAddDeviceHandler creates a new AddDeviceHandler
func NewAddDeviceHandler(deviceService *services.DeviceService) *AddDeviceHandler {
	return &AddDeviceHandler{deviceService: deviceService}
}

// ServeHTTP implements http.Handler interface
func (h *AddDeviceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Only allow POST method
	if r.Method != http.MethodPost {
		transport.SendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Parse request body
	var request models.AddDeviceRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Printf("Error decoding request: %v", err)
		transport.SendBadRequestError(w, "Invalid request format")
		return
	}

	// Validate request
	validationErrs := utils.Validate(request)
	if validationErrs != nil {
		log.Printf("Validation errors: %s", utils.ValidationErrorsToString(validationErrs))
		transport.SendBadRequestError(w, utils.CreateValidationError(validationErrs))
		return
	}

	// Get user ID from context (set by auth middleware)
	userID := middleware.GetUserID(r)
	if userID == "" {
		transport.SendError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	// Call service to add device
	if err := h.deviceService.AddDevice(r.Context(), request.DeviceID, request.DeviceName, userID); err != nil {
		log.Printf("Error adding device: %v", err)
		transport.SendError(w, http.StatusInternalServerError, "Failed to add device")
		return
	}

	transport.SendResponse(w, http.StatusCreated, "Device added successfully")
}

