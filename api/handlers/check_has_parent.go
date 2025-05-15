package handlers

import (
	"log"

	"github.com/gin-gonic/gin"

	"n1h41/zolaris-backend-app/internal/services"
	transport_gin "n1h41/zolaris-backend-app/internal/transport/gin"
)

type CheckHasParentIDHandler struct {
	UserService *services.UserService
}

func NewCheckHasParentIDHandler(userService *services.UserService) *CheckHasParentIDHandler {
	return &CheckHasParentIDHandler{UserService: userService}
}

// CheckHasParentIDHandler godoc
// @Summary Check if user has parent ID
// @Description Checks if the authenticated user has a parent ID set in their profile
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
		transport_gin.SendError(c, 500, "Failed to check parent ID")
		return
	}

	transport_gin.SendResponse(c, 200, gin.H{"has_parent_id": hasParentID})
}
