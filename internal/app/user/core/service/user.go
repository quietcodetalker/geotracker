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

// SetLocation sets user's location.
func (u userService) SetLocation(ctx context.Context, req port.SetLocationRequest) error {
	panic("implement me")
}
