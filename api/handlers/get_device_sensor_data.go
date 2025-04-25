package handlers

import (
	"encoding/json"
	"log"
	"n1h41/zolaris-backend-app/internal/models"
	"n1h41/zolaris-backend-app/internal/services"
	transport "n1h41/zolaris-backend-app/internal/transport/http"
	"n1h41/zolaris-backend-app/internal/utils"
	"net/http"
	"strconv"
)

// GetDeviceSensorDataHandler handles requests to get sensor data for a device
type GetDeviceSensorDataHandler struct {
	deviceService *services.DeviceService
}

// NewGetDeviceSensorDataHandler creates a new GetDeviceSensorDataHandler
func NewGetDeviceSensorDataHandler(deviceService *services.DeviceService) *GetDeviceSensorDataHandler {
	return &GetDeviceSensorDataHandler{deviceService: deviceService}
}

// ServeHTTP implements http.Handler interface
func (h *GetDeviceSensorDataHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Only allow POST method
	if r.Method != http.MethodPost {
		transport.SendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Parse request body
	var request models.GetDeviceSensorDataRequest
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

	// Parse the timestamp
	baseTimestamp, err := strconv.ParseInt(request.Timestamp, 10, 64)
	if err != nil {
		log.Printf("Error parsing timestamp: %v", err)
		transport.SendBadRequestError(w, "Invalid timestamp format")
		return
	}

	// Call service to get sensor data
	data, err := h.deviceService.GetDeviceSensorData(r.Context(), request.DeviceMacId, request.DateMode, baseTimestamp)
	if err != nil {
		log.Printf("Error getting sensor data: %v", err)
		transport.SendError(w, http.StatusInternalServerError, "Failed to retrieve sensor data")
		return
	}

	// Create response
	response := models.GetDeviceSensorDataResponse{
		Data: data,
	}

	transport.SendResponse(w, http.StatusOK, response)
}

