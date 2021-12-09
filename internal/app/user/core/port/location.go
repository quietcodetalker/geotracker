package port

import (
	"context"
	"gitlab.com/spacewalker/locations/internal/app/user/core/domain"
)

// GetUsersInRadiusRequest TODO: add description
// TODO: add validations
type GetUsersInRadiusRequest struct {
	Point  domain.Point
	Radius float64
}

// LocationService TODO: add description
type LocationService interface {
	GetUsersInRadius(ctx context.Context, req GetUsersInRadiusRequest) (GetUsersInRadiusResponse, error)
}

// GetUsersInRadiusResponse TODO: add description
type GetUsersInRadiusResponse struct {
}

// SetLocationArg is a param object of location repository SetLocation method.
type SetLocationArg struct {
	UserID int          `json:"user_id"`
	Point  domain.Point `json:"point"`
}

// UpdateLocationByUserIDArg is a param object of location repository UpdateLocationByUserID method.
type UpdateLocationByUserIDArg struct {
	UserID int          `json:"user_id"`
	Point  domain.Point `json:"point"`
}

// LocationRepository represents location repository.
type LocationRepository interface {
	SetLocation(ctx context.Context, arg SetLocationArg) (domain.Location, error)
	UpdateLocationByUserID(ctx context.Context, arg UpdateLocationByUserIDArg) (domain.Location, error)
}
