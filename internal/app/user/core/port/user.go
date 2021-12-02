package port

import (
	"context"
	"gitlab.com/spacewalker/locations/internal/app/user/core/domain"
)

// SetLocationRequest is a param object of user service SetLocation method.
type SetLocationRequest struct {
	Username  string  `json:"username" validate:"required,validusername"`
	Latitude  float64 `json:"latitude" validate:"required,gte=-180,lte=180"`
	Longitude float64 `json:"longitude" validate:"required,gte=-180,lte=180"`
}

// UserService represents user service.
type UserService interface {
	// SetLocation sets user's location by given username.
	SetLocation(ctx context.Context, req SetLocationRequest) error
}

// UserRepo represents user repository.
type UserRepo interface {
	// GetByUsername returns user with given username.
	// If user is not found returns ErrNotFound error.
	GetByUsername(ctx context.Context, username string) (*domain.User, error)
}
