package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"n1h41/zolaris-backend-app/internal/models"
	"n1h41/zolaris-backend-app/internal/services"
	transport_gin "n1h41/zolaris-backend-app/internal/transport/gin"
)

// GetCategoriesByTypeHandler handles requests to get categories by type
type GetCategoriesByTypeHandler struct {
	categoryService *services.CategoryService
}

// NewGetCategoriesByTypeHandler creates a new GetCategoriesByTypeHandler
func NewGetCategoriesByTypeHandler(categoryService *services.CategoryService) *GetCategoriesByTypeHandler {
	return &GetCategoriesByTypeHandler{categoryService: categoryService}
}

// HandleGin handles requests using Gin framework
// @Summary Get categories by type
// @Description Retrieve all categories of a specific type
// @Tags Category Management
// @Accept json
// @Produce json
// @Param type path string true "Category type"
// @Success 200 {array} models.CategoryResponse "List of categories"
// @Failure 500 {object} transport_gin.ErrorResponse "Internal server error"
// @Router /category/type/{type} [get]
func (h *GetCategoriesByTypeHandler) HandleGin(c *gin.Context) {
	// Get type from URL parameter
	categoryType := c.Param("type")
	if categoryType == "" {
		transport_gin.SendBadRequestError(c, "Category type is required")
		return
	}

	// Call service to get categories by type
	categories, err := h.categoryService.GetCategoriesByType(c.Request.Context(), categoryType)
	if err != nil {
		log.Printf("Error getting categories: %v", err)
		transport_gin.SendError(c, http.StatusInternalServerError, "Failed to get categories")
		return
	}

	// Return empty array if no categories found
	if categories == nil {
		categories = []models.CategoryResponse{}
	}

	c.JSON(http.StatusOK, categories)
}

