package handlers

import (
	"log"
	"n1h41/zolaris-backend-app/internal/models"
	"n1h41/zolaris-backend-app/internal/services"
	transport_gin "n1h41/zolaris-backend-app/internal/transport/gin"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ListAllCategoriesHandler struct {
	categoryService *services.CategoryService
}

func NewListAllCategoriesHandler(categoryService *services.CategoryService) *ListAllCategoriesHandler {
  return &ListAllCategoriesHandler{categoryService: categoryService}
}

// HandleGin handles requests using Gin framework
// @Summary Get all categories
// @Description Retrieve all categories
// @Tags Category Management
// @Produce json
// @Success 200 {array} models.CategoryResponse "List of categories"
// @Failure 500 {object} transport_gin.ErrorResponse "Internal server error"
// @Router /category/all [get]
func (h *ListAllCategoriesHandler) HandleGin(c *gin.Context) {
	// Call service to get categories by type
	categories, err := h.categoryService.GetAllCategories(c.Request.Context())
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
