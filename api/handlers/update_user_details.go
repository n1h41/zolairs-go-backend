package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"n1h41/zolaris-backend-app/internal/middleware"
	"n1h41/zolaris-backend-app/internal/models"
	"n1h41/zolaris-backend-app/internal/services"
	transport_gin "n1h41/zolaris-backend-app/internal/transport/gin"
	"n1h41/zolaris-backend-app/internal/utils"
)

// UpdateUserDetailsHandler represents a handler for updating user details
type UpdateUserDetailsHandler struct {
	UserService *services.UserService
}

// NewUpdateUserDetailsHandler creates a new update user details handler
func NewUpdateUserDetailsHandler(userService *services.UserService) *UpdateUserDetailsHandler {
	return &UpdateUserDetailsHandler{UserService: userService}
}

// UpdateUserDetailsHandler godoc
// @Summary Update user details
// @Description Updates or adds user details for the authenticated user
// @Accept json
// @Produce json
// @Param X-User-ID header string true "User ID"
// @Param request body models.UserDetailsRequest true "User details information"
// @Success 200 {object} map[string]interface{} "User details updated successfully"
// @Failure 400 {object} map[string]string "Error when request validation fails"
// @Failure 401 {object} map[string]string "Error when user is not authenticated"
// @Failure 500 {object} map[string]string "Error when updating user details fails"
// @Router /user/details [post]
func (h *UpdateUserDetailsHandler) HandleGin(c *gin.Context) {
	// Extract user ID from the context (set by auth middleware)
	userID := middleware.GetUserIDFromGin(c)
	if userID == "" {
		transport_gin.SendError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var req models.UserDetailsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		transport_gin.SendError(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate the request
	validationErrs := utils.Validate(req)
	if validationErrs != nil {
		log.Printf("Validation errors: %s", utils.ValidationErrorsToString(validationErrs))
		transport_gin.SendBadRequestError(c, utils.CreateValidationError(validationErrs))
		return
	}

	// Map request to user details model
	userDetails := models.UserDetails{
		City:      req.City,
		Country:   req.Country,
		Email:     req.Email,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Phone:     req.Phone,
		Region:    req.Region,
		Street1:   req.Street1,
		Street2:   req.Street2,
		Zip:       req.Zip,
	}

	// Update the user details
	err := h.UserService.UpdateUserDetails(c.Request.Context(), userID, &userDetails)
	if err != nil {
		log.Printf("Error updating user details: %v", err)
		transport_gin.SendError(c, http.StatusInternalServerError, "Failed to update user details")
		return
	}

	transport_gin.SendResponse(c, http.StatusOK, "User details updated successfully")
}
