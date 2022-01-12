package port

import (
	"context"
	"gitlab.com/spacewalker/locations/internal/app/location/core/domain"
)

// LocationServiceGetUsersInRadiusRequest TODO: add description
// TODO: add validations
type LocationServiceGetUsersInRadiusRequest struct {
	Point  domain.Point
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
	UserID int          `json:"user_id"`
	Point  domain.Point `json:"point"`
}

// LocationRepositoryUpdateLocationByUserIDRequest is a param object of location repository UpdateLocationByUserID method.
type LocationRepositoryUpdateLocationByUserIDRequest struct {
	UserID int          `json:"user_id"`
	Point  domain.Point `json:"point"`
}

// LocationRepository represents location repository.
type LocationRepository interface {
	SetLocation(ctx context.Context, arg LocationRepositorySetLocationRequest) (domain.Location, error)
	UpdateLocationByUserID(ctx context.Context, arg LocationRepositoryUpdateLocationByUserIDRequest) (domain.Location, error)
}
