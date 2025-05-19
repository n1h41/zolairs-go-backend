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

// GetUserDetailsHandler handles requests to get user details
type GetUserDetailsHandler struct {
	userService *services.UserService
}

// NewGetUserDetailsHandler creates a new GetUserDetailsHandler
func NewGetUserDetailsHandler(userService *services.UserService) *GetUserDetailsHandler {
	return &GetUserDetailsHandler{userService: userService}
}

// HandleGin handles GET /user/details requests
// @Summary Get user details
// @Description Retrieve authenticated user's profile information
// @Tags User Management
// @Accept json
// @Produce json
// @Param X-User-ID header string true "User ID"
// @Success 200 {object} dto.Response{data=dto.UserResponse} "User details retrieved successfully"
// @Failure 401 {object} dto.ErrorResponse "User not authenticated"
// @Failure 404 {object} dto.ErrorResponse "User not found"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Security ApiKeyAuth
// @Router /user/details [get]
func (h *GetUserDetailsHandler) HandleGin(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userID := middleware.GetUserIDFromGin(c)
	if userID == "" {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	// Get user from service
	user, err := h.userService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		log.Printf("Error retrieving user details: %v", err)
		response.InternalError(c, "Failed to retrieve user details")
		return
	}

	if user == nil {
		response.NotFound(c, "User not found")
		return
	}

	// Convert domain model to response DTO
	userResponse := mappers.UserToResponse(user)
	response.OK(c, userResponse, "User details retrieved successfully")
}

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

// CheckHasParentIDHandler handles requests to check if a user has a parent ID
type CheckHasParentIDHandler struct {
	UserService *services.UserService
}

// NewCheckHasParentIDHandler creates a new CheckHasParentIDHandler
func NewCheckHasParentIDHandler(userService *services.UserService) *CheckHasParentIDHandler {
	return &CheckHasParentIDHandler{UserService: userService}
}

// HandleGin handles GET /user/check-parent-id requests
// @Summary Check if user has parent ID
// @Description Checks if the authenticated user has a parent ID set in their profile
// @Tags User Management
// @Produce json
// @Param X-User-ID header string true "User ID"
// @Success 200 {object} map[string]bool "Returns has_parent_id flag"
// @Failure 400 {object} map[string]string "Error when user ID is not found in context"
// @Failure 500 {object} map[string]string "Error when checking parent ID fails"
// @Router /user/check-parent-id [get]
func (h *CheckHasParentIDHandler) HandleGin(c *gin.Context) {
	// Extract user ID from the request context
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(400, gin.H{"error": "User ID not found in context"})
		return
	}

	// Check if the user has a parent ID
	hasParentID, err := h.UserService.CheckHasParentID(c.Request.Context(), userID.(string))
	if err != nil {
		log.Printf("Error checking parent ID: %v", err)
		response.InternalError(c, "Failed to check parent ID")
		return
	}

	response.OK(c, gin.H{"has_parent_id": hasParentID}, "Success")
}
