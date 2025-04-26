package models

type DeviceResponse struct {
	DeviceID   string `json:"mac_address" dynamodbav:"mac_address"`
	UserID     string `json:"user_id" dynamodbav:"user_id"`
	DeviceName string `json:"device_name" dynamodbav:"device_name"`
}

type SensorData struct {
	Timestamp   int64  `json:"timestamp" dynamodbav:"timestamp"`
	Amperage    string `json:"amperage" dynamodbav:"amperage"`
	Temperature string `json:"temperature" dynamodbav:"temperature"`
	Humidity    string `json:"humidity" dynamodbav:"humidity"`
}

type GetDeviceSensorDataResponse struct {
	Data []SensorData `json:"data" dynamodbav:"data"`
}
