package models

type UserDetails struct {
	City      string `json:"city" dynamodbav:"city"`
	Country   string `json:"country" dynamodbav:"country"`
	Email     string `json:"email" dynamodbav:"email"`
	FirstName string `json:"firstName" dynamodbav:"firstName"`
	LastName  string `json:"lastName" dynamodbav:"lastName"`
	Phone     string `json:"phone" dynamodbav:"phone"`
	Region    string `json:"region" dynamodbav:"region"`
	Street1   string `json:"street1" dynamodbav:"street1"`
	Street2   string `json:"street2" dynamodbav:"street2"`
	Zip       string `json:"zip" dynamodbav:"zip"`
}

type UserResponse struct {
	UserID      string      `json:"user_id" dynamodbav:"user_id"`
	UserDetails UserDetails `json:"userDetails" dynamodbav:"userDetails"`
}

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

type CategoryResponse struct {
	Name string `json:"name" dynamodbav:"name"`
	Type string `json:"type" dynamodbav:"type"`
}


