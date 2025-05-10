package services

import (
	"context"
	"log"
	"n1h41/zolaris-backend-app/internal/models"
	"n1h41/zolaris-backend-app/internal/repositories"
	"time"
)

// DeviceService handles business logic for device operations
type DeviceService struct {
	deviceRepo *repositories.DeviceRepository
}

// NewDeviceService creates a new device service instance
func NewDeviceService(deviceRepo *repositories.DeviceRepository) *DeviceService {
	return &DeviceService{deviceRepo: deviceRepo}
}

// AddDevice handles the business logic for adding a new device
func (s *DeviceService) AddDevice(ctx context.Context, deviceID, deviceName, userID string) error {
	// Add any business logic here (validation, etc.)
	log.Printf("Adding device %s for user %s", deviceID, userID)
	return s.deviceRepo.AddDevice(ctx, deviceID, deviceName, userID)
}

// GetUserDevices retrieves all devices for a user
func (s *DeviceService) GetUserDevices(ctx context.Context, userID string) ([]models.DeviceResponse, error) {
	log.Printf("Getting devices for user %s", userID)
	return s.deviceRepo.GetDevicesByUserID(ctx, userID)
}

// GetDeviceSensorData retrieves sensor data for a device within a time range
func (s *DeviceService) GetDeviceSensorData(ctx context.Context, macID, dateMode string, timestamp int64) ([]models.SensorData, error) {
	// Calculate time range based on dateMode
	startTime, endTime := s.calculateTimeRange(timestamp, dateMode)
	log.Printf("Getting sensor data for device %s from %d to %d", macID, startTime, endTime)
	return s.deviceRepo.GetSensorData(ctx, macID, startTime, endTime)
}

// calculateTimeRange calculates a time range looking backward from the provided timestamp
func (s *DeviceService) calculateTimeRange(baseTimeMs int64, dateMode string) (int64, int64) {
	// Convert milliseconds to seconds and nanoseconds for time package
	seconds := baseTimeMs / 1000
	nanoseconds := (baseTimeMs % 1000) * 1000000

	endTime := time.Unix(seconds, nanoseconds).UTC() // The provided timestamp becomes the end time
	var startTime time.Time

	switch dateMode {
	case "hourly":
		// Look back 1 hour from the provided timestamp
		startTime = endTime.Add(-1 * time.Hour)
	case "daily":
		// Look back 24 hours from the provided timestamp
		startTime = endTime.Add(-24 * time.Hour)
	case "weekly":
		// Look back 7 days from the provided timestamp
		startTime = endTime.Add(-7 * 24 * time.Hour)
	case "monthly":
		// Look back approximately 30 days from the provided timestamp
		startTime = endTime.AddDate(0, -1, 0)
	case "yearly":
		// Look back 1 year from the provided timestamp
		startTime = endTime.AddDate(-1, 0, 0)
	default:
		// Default to daily if unrecognized mode
		startTime = endTime.Add(-24 * time.Hour)
	}

	// Return timestamps in milliseconds
	return startTime.UnixMilli(), endTime.UnixMilli()
}
