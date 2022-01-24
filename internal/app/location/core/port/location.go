package port

import (
	"context"
	"gitlab.com/spacewalker/locations/internal/app/location/core/domain"
	"gitlab.com/spacewalker/locations/internal/pkg/geo"
)

// LocationServiceGetUsersInRadiusRequest TODO: add description
// TODO: add validations
type LocationServiceGetUsersInRadiusRequest struct {
	Point  geo.Point
	Radius float64
}

// LocationServiceGetUsersInRadiusResponse TODO: add description
type LocationServiceGetUsersInRadiusResponse struct {
}

// LocationService TODO: add description
type LocationService interface {
	GetUsersInRadius(ctx context.Context, req LocationServiceGetUsersInRadiusRequest) (LocationServiceGetUsersInRadiusResponse, error)
}

// LocationRepositorySetLocationRequest is a param object of location repository SetLocation method.
type LocationRepositorySetLocationRequest struct {
	UserID int       `json:"user_id"`
	Point  geo.Point `json:"point"`
}

// LocationRepositoryUpdateLocationByUserIDRequest is a param object of location repository UpdateLocationByUserID method.
type LocationRepositoryUpdateLocationByUserIDRequest struct {
	UserID int       `json:"user_id"`
	Point  geo.Point `json:"point"`
}

// LocationRepository represents location repository.
type LocationRepository interface {
	GetLocation(ctx context.Context, userID int) (domain.Location, error)
	SetLocation(ctx context.Context, arg LocationRepositorySetLocationRequest) (domain.Location, error)
	UpdateLocationByUserID(ctx context.Context, arg LocationRepositoryUpdateLocationByUserIDRequest) (domain.Location, error)
}
