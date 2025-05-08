package services

import (
	"context"
	"errors"
	"log"
	"n1h41/zolaris-backend-app/internal/models"
	"n1h41/zolaris-backend-app/internal/repositories"
)

// CategoryService handles business logic for category operations
type CategoryService struct {
	categoryRepo *repositories.CategoryRepository
}

// NewCategoryService creates a new category service instance
func NewCategoryService(categoryRepo *repositories.CategoryRepository) *CategoryService {
	return &CategoryService{categoryRepo: categoryRepo}
}

// AddCategory handles the business logic for adding a new category
func (s *CategoryService) AddCategory(ctx context.Context, name, categoryType string) error {
	log.Printf("Adding category %s of type %s", name, categoryType)

	// Check if category already exists
	existingCategory, err := s.categoryRepo.GetCategoryByName(ctx, name)
	if err != nil {
		return err
	}

	if existingCategory != nil {
		return errors.New("category with this name already exists")
	}

	return s.categoryRepo.AddCategory(ctx, name, categoryType)
}

// GetCategoryByName retrieves a category by its name
func (s *CategoryService) GetCategoryByName(ctx context.Context, name string) (*models.CategoryResponse, error) {
	log.Printf("Getting category with name %s", name)
	return s.categoryRepo.GetCategoryByName(ctx, name)
}

// GetCategoriesByType retrieves all categories of a specific type
func (s *CategoryService) GetCategoriesByType(ctx context.Context, categoryType string) ([]models.CategoryResponse, error) {
	log.Printf("Getting categories of type %s", categoryType)
	return s.categoryRepo.GetCategoriesByType(ctx, categoryType)
}

func (s *CategoryService) GetAllCategories(ctx context.Context) ([]models.CategoryResponse, error) {
	log.Println("List all categories")
	return s.categoryRepo.ListAllCategories(ctx)
}
