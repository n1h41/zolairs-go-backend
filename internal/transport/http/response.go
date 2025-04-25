package http

import (
	"encoding/json"
	"log"
	"net/http"
)

// SendResponse sends a successful JSON response
func SendResponse(w http.ResponseWriter, statusCode int, data any) {
	resp := map[string]any{
		"status": true,
	}

	// Handle different types of data
	switch data.(type) {
	case string:
		resp["message"] = data.(string)
	case any:
		resp["data"] = data
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

// SendError sends an error JSON response
func SendError(w http.ResponseWriter, statusCode int, err any) {
	log.Printf("Error: %v", err)

	errorResp := map[string]any{
		"status": false,
	}

	// Handle different types of errors
	switch err.(type) {
	case error:
		errorResp["message"] = err.(error).Error()
	case string:
		errorResp["message"] = err.(string)
	default:
		errorResp["message"] = "An unknown error occurred"
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if encodeErr := json.NewEncoder(w).Encode(errorResp); encodeErr != nil {
		log.Printf("Error encoding error response: %v", encodeErr)
		http.Error(w, "Error encoding error response", http.StatusInternalServerError)
	}
}

// SendBadRequestError sends a 400 Bad Request error response
func SendBadRequestError(w http.ResponseWriter, err any) {
	SendError(w, http.StatusBadRequest, err)
}

