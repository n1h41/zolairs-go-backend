package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"n1h41/zolaris-backend-app/internal/services"
	"n1h41/zolaris-backend-app/internal/transport/dto"
	"n1h41/zolaris-backend-app/internal/transport/response"
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
// @Param category body dto.CategoryRequest true "Category information"
// @Success 201 {object} dto.Response "Category added successfully"
// @Failure 400 {object} dto.ErrorResponse "Validation error"
// @Failure 409 {object} dto.ErrorResponse "Category already exists"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /category/add [post]
func (h *AddCategoryHandler) HandleGin(c *gin.Context) {
	// Parse request body
	var request dto.CategoryRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		log.Printf("Error decoding request: %v", err)
		response.BadRequest(c, "Invalid request format")
		return
	}

	// Validate request
	validationErrs := utils.Validate(request)
	if validationErrs != nil {
		log.Printf("Validation errors: %s", utils.ValidationErrorsToString(validationErrs))
		response.ValidationErrors(c, utils.CreateDtoValidationErrors(validationErrs))
		return
	}

	// Call service to add category
	if err := h.categoryService.AddCategory(c.Request.Context(), request.Name, request.Type); err != nil {
		if err.Error() == "category with this name already exists" {
			response.Error(c, http.StatusConflict, "Category with this name already exists", "CONFLICT")
			return
		}
		log.Printf("Error adding category: %v", err)
		response.InternalError(c, "Failed to add category")
		return
	}

	response.Created(c, nil, "Category added successfully")
}
