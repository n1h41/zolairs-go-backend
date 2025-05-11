package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"n1h41/zolaris-backend-app/internal/middleware"
	"n1h41/zolaris-backend-app/internal/services"
	transport_gin "n1h41/zolaris-backend-app/internal/transport/gin"
)

// GetUserDetailsHandler represents a handler for retrieving user details
type GetUserDetailsHandler struct {
	UserService *services.UserService
}

// NewGetUserDetailsHandler creates a new get user details handler
func NewGetUserDetailsHandler(userService *services.UserService) *GetUserDetailsHandler {
	return &GetUserDetailsHandler{UserService: userService}
}

// GetUserDetailsHandler godoc
// @Summary Get user details
// @Description Retrieves the user details for the authenticated user
// @Produce json
// @Param X-User-ID header string true "User ID"
// @Success 200 {object} models.UserDetails "User details retrieved successfully"
// @Success 200 {object} map[string]interface{} "No user details found (null response)"
// @Failure 401 {object} map[string]string "Error when user is not authenticated"
// @Failure 500 {object} map[string]string "Error when retrieving user details fails"
// @Router /user/details [get]
func (h *GetUserDetailsHandler) HandleGin(c *gin.Context) {
	// Extract user ID from the context (set by auth middleware)
	userID := middleware.GetUserIDFromGin(c)
	if userID == "" {
		transport_gin.SendError(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	// Get the user details
	userDetails, err := h.UserService.GetUserDetails(c.Request.Context(), userID)
	if err != nil {
		log.Printf("Error retrieving user details: %v", err)
		transport_gin.SendError(c, http.StatusInternalServerError, "Failed to retrieve user details")
		return
	}

	if userDetails == nil {
		// User has no details yet
		transport_gin.SendResponse(c, http.StatusOK, "No user details found")
		return
	}

	transport_gin.SendResponse(c, http.StatusOK, gin.H{"userDetails": userDetails})
}
