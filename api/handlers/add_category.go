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

// AddCategoryHandler handles requests to add a new category
type AddCategoryHandler struct {
	categoryService *services.CategoryService
}

// NewAddCategoryHandler creates a new AddCategoryHandler
func NewAddCategoryHandler(categoryService *services.CategoryService) *AddCategoryHandler {
	return &AddCategoryHandler{categoryService: categoryService}
}

// HandleGin handles requests using Gin framework
// @Summary Add a new category
// @Description Register a new category
// @Tags Category Management
// @Accept json
// @Produce json
// @Param category body models.AddCategoryRequest true "Category information"
// @Success 201 {object} transport_gin.Response "Category added successfully"
// @Failure 400 {object} transport_gin.ErrorResponse "Validation error"
// @Failure 409 {object} transport_gin.ErrorResponse "Category already exists"
// @Failure 500 {object} transport_gin.ErrorResponse "Internal server error"
// @Router /category/add [post]
func (h *AddCategoryHandler) HandleGin(c *gin.Context) {
	// Parse request body
	var request models.AddCategoryRequest
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

	// Call service to add category
	if err := h.categoryService.AddCategory(c.Request.Context(), request.Name, request.Type); err != nil {
		if err.Error() == "category with this name already exists" {
			transport_gin.SendError(c, http.StatusConflict, "Category with this name already exists")
			return
		}
		log.Printf("Error adding category: %v", err)
		transport_gin.SendError(c, http.StatusInternalServerError, "Failed to add category")
		return
	}

	transport_gin.SendResponse(c, http.StatusCreated, "Category added successfully")
}
