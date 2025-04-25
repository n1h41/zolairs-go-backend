package models

type DeviceResponse struct {
	DeviceID   string `json:"deviceId"`
	UserID     string `json:"userId"`
	DeviceName string `json:"deviceName"`
}

type SensorData struct {
	Timestamp   int64    `json:"timestamp"`
	Amperage    string `json:"amperage"`
	Temperature string `json:"temperature"`
	Humidity    string `json:"humidity"`
}

type GetDeviceSensorDataResponse struct {
	Data []SensorData `json:"data"`
}
