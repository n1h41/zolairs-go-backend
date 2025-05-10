package services

import "n1h41/zolaris-backend-app/internal/repositories"

type UserService struct {
  deviceRepo *repositories.DeviceRepository
}

func NewUserService(deviceRepo *repositories.DeviceRepository) *UserService {
  return &UserService{deviceRepo: deviceRepo}
}
