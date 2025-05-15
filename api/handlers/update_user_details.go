package handlers

import (
	"log"

	"github.com/gin-gonic/gin"

	"n1h41/zolaris-backend-app/internal/middleware"
	"n1h41/zolaris-backend-app/internal/services"
	"n1h41/zolaris-backend-app/internal/transport/dto"
	"n1h41/zolaris-backend-app/internal/transport/mappers"
	"n1h41/zolaris-backend-app/internal/transport/response"
	"n1h41/zolaris-backend-app/internal/utils"
)

// UpdateUserDetailsHandler handles requests to update user details
type UpdateUserDetailsHandler struct {
	userService *services.UserService
}

// NewUpdateUserDetailsHandler creates a new UpdateUserDetailsHandler
func NewUpdateUserDetailsHandler(userService *services.UserService) *UpdateUserDetailsHandler {
	return &UpdateUserDetailsHandler{userService: userService}
}

// HandleGin handles POST /user/details requests
// @Summary Update user details
// @Description Update the authenticated user's profile information
// @Tags User Management
// @Accept json
// @Produce json
// @Param X-User-ID header string true "User ID"
// @Param user body dto.UserDetailsRequest true "User details"
// @Success 200 {object} dto.Response{data=dto.UserResponse} "User details updated successfully"
// @Failure 400 {object} dto.ErrorResponse "Validation error"
// @Failure 401 {object} dto.ErrorResponse "User not authenticated"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Security ApiKeyAuth
// @Router /user/details [post]
func (h *UpdateUserDetailsHandler) HandleGin(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID := middleware.GetUserIDFromGin(c)
	if userID == "" {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	// Parse request body
	var request dto.UserDetailsRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Printf("Error decoding request: %v", err)
		response.BadRequest(c, "Invalid request format")
		return
	}

	// Validate request
	validationErrs := utils.Validate(request)
	if validationErrs != nil {
		log.Printf("Validation errors: %s", utils.ValidationErrorsToString(validationErrs))

		// Convert validation errors to DTO format
		var validationErrDTOs []dto.ValidationError
		for _, item := range validationErrs {
			validationErrDTOs = append(validationErrDTOs, dto.ValidationError{
				Field:   item.Field,
				Message: item.Error,
			})
		}

		response.ValidationErrors(c, validationErrDTOs)
		return
	}

	// Update user details
	updatedUser, err := h.userService.UpdateUserDetails(c.Request.Context(), userID, &request)
	if err != nil {
		log.Printf("Error updating user details: %v", err)
		response.InternalError(c, "Failed to update user details")
		return
	}

	// Convert domain model to response DTO
	userResponse := mappers.UserToResponse(updatedUser)
	response.OK(c, userResponse, "User details updated successfully")
}
