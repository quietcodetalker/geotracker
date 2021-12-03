package service

import (
	"context"
	"gitlab.com/spacewalker/locations/internal/app/user/core/port"
)

type userService struct{}

// NewUserService creates instance of UserService and returns it's pointer.
func NewUserService() port.UserService {
	return &userService{}
}

// SetUserLocation sets user's location.
func (s *userService) SetUserLocation(ctx context.Context, req port.SetUserLocationRequest) (port.SetUserLocationResponse, error) {
	panic("implement me")
}
