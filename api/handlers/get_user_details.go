package handlers

import (
	"log"

	"github.com/gin-gonic/gin"
	"n1h41/zolaris-backend-app/internal/middleware"
	"n1h41/zolaris-backend-app/internal/services"
	"n1h41/zolaris-backend-app/internal/transport/mappers"
	"n1h41/zolaris-backend-app/internal/transport/response"
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

