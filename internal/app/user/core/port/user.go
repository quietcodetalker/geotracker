package port

import (
	"context"
	"gitlab.com/spacewalker/locations/internal/app/user/core/domain"
)

// SetUserLocationRequest is a param object of user service SetUserLocation method.
type SetUserLocationRequest struct {
	Username  string  `json:"username" validate:"required,validusername"`
	Latitude  float64 `json:"latitude" validate:"required,gte=-180,lte=180"`
	Longitude float64 `json:"longitude" validate:"required,gte=-180,lte=180"`
}

// SetUserLocationResponse represents response from user service SetUserLocation method.
type SetUserLocationResponse struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// UserService represents user service.
type UserService interface {
	SetUserLocation(ctx context.Context, req SetUserLocationRequest) (SetUserLocationResponse, error)
}

// CreateUserArg is a param object of use repository CreateUser method.
type CreateUserArg struct {
	Username string `json:"username"`
}

// SetUserLocationArg is a param object of user repository SetUserLocation method.
type SetUserLocationArg struct {
	Username string       `json:"username"`
	Point    domain.Point `json:"point"`
}

// UserRepository represents user repository.
type UserRepository interface {
	CreateUser(ctx context.Context, arg CreateUserArg) (domain.User, error)
	GetByUsername(ctx context.Context, username string) (domain.User, error)
	SetUserLocation(ctx context.Context, arg SetUserLocationArg) (domain.Location, error)
}
