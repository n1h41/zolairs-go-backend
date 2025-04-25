package handlers

import (
	"encoding/json"
	"log"
	"n1h41/zolaris-backend-app/internal/models"
	"n1h41/zolaris-backend-app/internal/services"
	transport "n1h41/zolaris-backend-app/internal/transport/http"
	"n1h41/zolaris-backend-app/internal/utils"
	"net/http"
)

// AttachIotPolicyHandler handles requests to attach an IoT policy to an identity
type AttachIotPolicyHandler struct {
	policyService *services.PolicyService
}

// NewAttachIotPolicyHandler creates a new AttachIotPolicyHandler
func NewAttachIotPolicyHandler(policyService *services.PolicyService) *AttachIotPolicyHandler {
	return &AttachIotPolicyHandler{policyService: policyService}
}

// ServeHTTP implements http.Handler interface
func (h *AttachIotPolicyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Only allow POST method
	if r.Method != http.MethodPost {
		transport.SendError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	// Parse request body
	var request models.AttachIotPolicyRequest
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

	// Call service to attach policy
	if err := h.policyService.AttachIoTPolicy(r.Context(), request.IdentityId); err != nil {
		log.Printf("Error attaching IoT policy: %v", err)
		transport.SendError(w, http.StatusInternalServerError, "Failed to attach IoT policy")
		return
	}

	transport.SendResponse(w, http.StatusOK, "IoT policy attached successfully")
}

