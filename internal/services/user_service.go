package services

import (
	"context"
	"log"
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
