package handlers

import (
	"log"

	"github.com/gin-gonic/gin"

	"n1h41/zolaris-backend-app/internal/services"
	"n1h41/zolaris-backend-app/internal/transport/dto"
	"n1h41/zolaris-backend-app/internal/transport/response"
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
// @Param request body dto.SensorDataRequest true "Request parameters"
// @Success 200 {object} dto.Response{data=[]dto.SensorDataResponse} "Sensor data for the device"
// @Failure 400 {object} dto.ErrorResponse "Invalid request or validation error"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /device/sensor-data [post]
func (h *GetDeviceSensorDataHandler) HandleGin(c *gin.Context) {
	// Parse request body
	var request dto.SensorDataRequest
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

	// Use timestamp directly from request
	baseTimestamp := request.Timestamp

	// Call service to get sensor data
	data, err := h.deviceService.GetDeviceSensorData(c.Request.Context(), request.DeviceMacID, request.DateMode, baseTimestamp)
	if err != nil {
		log.Printf("Error getting sensor data: %v", err)
		response.InternalError(c, "Failed to retrieve sensor data")
		return
	}

	// Use the data directly in the response
	response.OK(c, data, "Data retrieved successfully")
}
