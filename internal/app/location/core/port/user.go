//go:generate mockgen -destination=mock/mock_user.go -package=mock . UserRepository,UserService

package port

import (
	"context"
	"gitlab.com/spacewalker/locations/internal/app/location/core/domain"
)

// UserServiceSetUserLocationRequest is a param object of user service SetUserLocation method.
type UserServiceSetUserLocationRequest struct {
	Username  string  `json:"username" validate:"required,validusername"`
	Latitude  float64 `json:"latitude" validate:"validlatitude"`
	Longitude float64 `json:"longitude" validate:"validlongitude"`
}

// UserServiceSetUserLocationResponse represents response from user service SetUserLocation method.
type UserServiceSetUserLocationResponse struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// UserServiceListUsersInRadiusRequest TODO: add description
type UserServiceListUsersInRadiusRequest struct {
	Point     domain.Point `json:"point" validate:"validgeopoint"`
	Radius    float64      `json:"radius" validate:""`
	PageToken string       `json:"page_token" validate:"required_without=PageSize"`
	PageSize  int          `json:"page_size" validate:"required_without=PageToken"`
}

// UserServiceListUsersInRadiusResponse TODO: add description
type UserServiceListUsersInRadiusResponse struct {
	Users         []domain.User
	NextPageToken string
}

// UserService represents user service.
type UserService interface {
	SetUserLocation(ctx context.Context, req UserServiceSetUserLocationRequest) (UserServiceSetUserLocationResponse, error)
	ListUsersInRadius(ctx context.Context, req UserServiceListUsersInRadiusRequest) (UserServiceListUsersInRadiusResponse, error)
}

// CreateUserArg is a param object of use repository CreateUser method.
type CreateUserArg struct {
	Username string `json:"username"`
}

// UserRepositorySetUserLocationRequest is a param object of user repository SetUserLocation method.
type UserRepositorySetUserLocationRequest struct {
	Username string       `json:"username"`
	Point    domain.Point `json:"point"`
}

// UserRepositoryListUsersInRadiusRequest TODO: add description
type UserRepositoryListUsersInRadiusRequest struct {
	Point     domain.Point
	Radius    float64
	PageToken int
	PageSize  int
}

// UserRepositoryListUsersInRadiusResponse TODO: add description
type UserRepositoryListUsersInRadiusResponse struct {
	Users         []domain.User
	NextPageToken int
}

// UserRepository represents user repository.
type UserRepository interface {
	CreateUser(ctx context.Context, arg CreateUserArg) (domain.User, error)
	GetByUsername(ctx context.Context, username string) (domain.User, error)
	SetUserLocation(ctx context.Context, arg UserRepositorySetUserLocationRequest) (domain.Location, error)
	ListUsersInRadius(ctx context.Context, arg UserRepositoryListUsersInRadiusRequest) (UserRepositoryListUsersInRadiusResponse, error)
}
