package services

import (
	"context"
	"fmt"
	"log"

	"n1h41/zolaris-backend-app/internal/models"
	"n1h41/zolaris-backend-app/internal/repositories"
)

type UserService struct {
	userRepo *repositories.UserRepository
}

func NewUserService(userRepo *repositories.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

// HasParentID checks if a user has a parent ID (for hierarchical user relationships)
func (s *UserService) HasParentID(ctx context.Context, userID string) (bool, error) {
	log.Printf("Checking if user %s has a parent ID", userID)
	return s.userRepo.HasParentID(ctx, userID)
}

// UpdateUserDetails adds or updates user details for a specific user
func (s *UserService) UpdateUserDetails(ctx context.Context, userID string, details *models.UserDetails) error {
	log.Printf("Updating details for user %s", userID)

	// Validate required user detail fields
	if details.Email == "" || details.FirstName == "" || details.LastName == "" {
		return fmt.Errorf("missing required user detail fields")
	}

	err := s.userRepo.UpdateUserDetails(ctx, userID, details)
	if err != nil {
		log.Printf("Error in repository while updating user details: %v", err)
		return err
	}

	return nil
}

// GetUserDetails retrieves user details for a specific user
func (s *UserService) GetUserDetails(ctx context.Context, userID string) (*models.UserDetails, error) {
	log.Printf("Getting details for user %s", userID)

	if userID == "" {
		return nil, fmt.Errorf("userID cannot be empty")
	}

	userDetails, err := s.userRepo.GetUserDetails(ctx, userID)
	if err != nil {
		log.Printf("Error retrieving user details: %v", err)
		return nil, fmt.Errorf("failed to retrieve user details: %w", err)
	}

	return userDetails, nil
}
