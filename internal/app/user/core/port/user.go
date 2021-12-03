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

// SetUserLocationResponse represents response from user service SetUserLocation method.
type SetUserLocationResponse struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// UserService represents user service.
type UserService interface {
	// SetUserLocation sets user's location by given username.
	SetUserLocation(ctx context.Context, req SetUserLocationRequest) (SetUserLocationResponse, error)
}

// SetUserLocationArg is a param object of user repository SetUserLocation method.
type SetUserLocationArg struct {
	Username  string  `json:"username"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// UserRepo represents user repository.
type UserRepo interface {
	// GetByUsername returns user with given username.
	// If user is not found returns ErrNotFound error.
	GetByUsername(ctx context.Context, username string) (domain.User, error)

	// SetUserLocation gets User by given username and updates Location by user ID
	//with provided coordinates within a single database transaction.
	SetUserLocation(ctx context.Context, arg SetUserLocationArg) (domain.Location, error)
}
