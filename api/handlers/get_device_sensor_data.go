package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"n1h41/zolaris-backend-app/internal/models"
	"n1h41/zolaris-backend-app/internal/services"
	transport_gin "n1h41/zolaris-backend-app/internal/transport/gin"
	"n1h41/zolaris-backend-app/internal/utils"
)

// GetDeviceSensorDataHandler handles requests to get sensor data for a device
type GetDeviceSensorDataHandler struct {
	deviceService *services.DeviceService
}

// NewGetDeviceSensorDataHandler creates a new GetDeviceSensorDataHandler
func NewGetDeviceSensorDataHandler(deviceService *services.DeviceService) *GetDeviceSensorDataHandler {
	return &GetDeviceSensorDataHandler{deviceService: deviceService}
}

// HandleGin handles requests using Gin framework
// @Summary Get device sensor data
// @Description Retrieve sensor data for a specific device with time filtering
// @Tags Device Data
// @Accept json
// @Produce json
// @Param request body models.GetDeviceSensorDataRequest true "Request parameters"
// @Success 200 {object} models.GetDeviceSensorDataResponse "Sensor data for the device"
// @Failure 400 {object} transport_gin.ErrorResponse "Invalid request or validation error"
// @Failure 500 {object} transport_gin.ErrorResponse "Internal server error"
// @Router /device/sensor-data [post]
func (h *GetDeviceSensorDataHandler) HandleGin(c *gin.Context) {
	// Parse request body
	var request models.GetDeviceSensorDataRequest
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

	// Parse the timestamp
	baseTimestamp, err := strconv.ParseInt(request.Timestamp, 10, 64)
	if err != nil {
		log.Printf("Error parsing timestamp: %v", err)
		transport_gin.SendBadRequestError(c, "Invalid timestamp format")
		return
	}

	// Call service to get sensor data
	data, err := h.deviceService.GetDeviceSensorData(c.Request.Context(), request.DeviceMacId, request.DateMode, baseTimestamp)
	if err != nil {
		log.Printf("Error getting sensor data: %v", err)
		transport_gin.SendError(c, http.StatusInternalServerError, "Failed to retrieve sensor data")
		return
	}

	// Create response
	response := models.GetDeviceSensorDataResponse{
		Data: data,
	}

	transport_gin.SendResponse(c, http.StatusOK, response)
}
