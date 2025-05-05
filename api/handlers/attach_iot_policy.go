package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"n1h41/zolaris-backend-app/internal/models"
	"n1h41/zolaris-backend-app/internal/services"
	transport_gin "n1h41/zolaris-backend-app/internal/transport/gin"
	transport_http "n1h41/zolaris-backend-app/internal/transport/http"
	"n1h41/zolaris-backend-app/internal/utils"
)

// AttachIotPolicyHandler handles requests to attach an IoT policy to an identity
type AttachIotPolicyHandler struct {
	policyService *services.PolicyService
}

// NewAttachIotPolicyHandler creates a new AttachIotPolicyHandler
func NewAttachIotPolicyHandler(policyService *services.PolicyService) *AttachIotPolicyHandler {
	return &AttachIotPolicyHandler{policyService: policyService}
}

// ServeHTTP implements http.Handler interface (for backward compatibility)
func (h *AttachIotPolicyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Only allow POST method
	if r.Method != http.MethodPost {
		transport_http.SendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Parse request body
	var request models.AttachIotPolicyRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Printf("Error decoding request: %v", err)
		transport_http.SendBadRequestError(w, "Invalid request format")
		return
	}

	// Validate request
	validationErrs := utils.Validate(request)
	if validationErrs != nil {
		log.Printf("Validation errors: %s", utils.ValidationErrorsToString(validationErrs))
		transport_http.SendBadRequestError(w, utils.CreateValidationError(validationErrs))
		return
	}

	// Call service to attach policy
	if err := h.policyService.AttachIoTPolicy(r.Context(), request.IdentityId); err != nil {
		log.Printf("Error attaching IoT policy: %v", err)
		transport_http.SendError(w, http.StatusInternalServerError, "Failed to attach IoT policy")
		return
	}

	transport_http.SendResponse(w, http.StatusOK, "IoT policy attached successfully")
}

// HandleGin handles requests using Gin framework
func (h *AttachIotPolicyHandler) HandleGin(c *gin.Context) {
	// Parse request body
	var request models.AttachIotPolicyRequest
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

	// Call service to attach policy
	if err := h.policyService.AttachIoTPolicy(c.Request.Context(), request.IdentityId); err != nil {
		log.Printf("Error attaching IoT policy: %v", err)
		transport_gin.SendError(c, http.StatusInternalServerError, "Failed to attach IoT policy")
		return
	}

	transport_gin.SendResponse(c, http.StatusOK, "IoT policy attached successfully")
}
