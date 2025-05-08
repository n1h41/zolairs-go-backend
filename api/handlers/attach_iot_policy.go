package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"n1h41/zolaris-backend-app/internal/models"
	"n1h41/zolaris-backend-app/internal/services"
	transport_gin "n1h41/zolaris-backend-app/internal/transport/gin"
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

// HandleGin handles requests using Gin framework
// @Summary Attach IoT policy
// @Description Attach an AWS IoT policy to a Cognito identity
// @Tags Policy Management
// @Accept json
// @Produce json
// @Param request body models.AttachIotPolicyRequest true "Identity information"
// @Success 200 {object} transport_gin.Response "IoT policy attached successfully"
// @Failure 400 {object} transport_gin.ErrorResponse "Invalid request or validation error"
// @Failure 500 {object} transport_gin.ErrorResponse "Failed to attach IoT policy"
// @Router /device/attach-policy [post]
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
