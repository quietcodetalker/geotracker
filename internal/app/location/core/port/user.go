//go:generate mockgen -destination=mock/mock_user.go -package=mock . UserRepository,UserService

package port

import (
	"context"

	"gitlab.com/spacewalker/geotracker/internal/app/location/core/domain"
	"gitlab.com/spacewalker/geotracker/internal/pkg/geo"
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
	Point     geo.Point `json:"point" validate:"validgeopoint"`
	Radius    float64   `json:"radius" validate:"gte=0"`
	PageToken string    `json:"page_token" validate:"required_without=PageSize"`
	PageSize  int       `json:"page_size" validate:"required_without=PageToken"`
}

// UserServiceListUsersInRadiusResponse TODO: add description
type UserServiceListUsersInRadiusResponse struct {
	Users         []domain.User `json:"users"`
	NextPageToken string        `json:"next_page_token"`
}

// UserService represents user service.
type UserService interface {
	GetByUsername(ctx context.Context, username string) (domain.User, error)
	SetUserLocation(ctx context.Context, req UserServiceSetUserLocationRequest) (UserServiceSetUserLocationResponse, error)
	ListUsersInRadius(ctx context.Context, req UserServiceListUsersInRadiusRequest) (UserServiceListUsersInRadiusResponse, error)
}

// CreateUserArg is a param object of use repository CreateUser method.
type CreateUserArg struct {
	Username string `json:"username"`
}

// UserRepositorySetUserLocationRequest is a param object of user repository SetUserLocation method.
type UserRepositorySetUserLocationRequest struct {
	Username string    `json:"username"`
	Point    geo.Point `json:"point"`
}

// UserRepositoryListUsersInRadiusRequest TODO: add description
type UserRepositoryListUsersInRadiusRequest struct {
	Point     geo.Point
	Radius    float64
	PageToken int
	PageSize  int
}

// UserRepositoryListUsersInRadiusResponse TODO: add description
type UserRepositoryListUsersInRadiusResponse struct {
	Users         []domain.User
	NextPageToken int
}

// UserRepositorySetUserLocationResponse TODO: add description
type UserRepositorySetUserLocationResponse struct {
	User         domain.User
	PrevLocation domain.Location
	Location     domain.Location
}

// UserRepository represents user repository.
type UserRepository interface {
	CreateUser(ctx context.Context, arg CreateUserArg) (domain.User, error)
	GetByUsername(ctx context.Context, username string) (domain.User, error)
	SetUserLocation(ctx context.Context, arg UserRepositorySetUserLocationRequest) (UserRepositorySetUserLocationResponse, error)
	ListUsersInRadius(ctx context.Context, arg UserRepositoryListUsersInRadiusRequest) (UserRepositoryListUsersInRadiusResponse, error)
}
