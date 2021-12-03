package port

import (
	"context"
	"gitlab.com/spacewalker/locations/internal/app/user/core/domain"
)

// SetLocationArg is a param object of user repository SetLocation method.
type SetLocationArg struct {
	UserID    int     `json:"user_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// LocationRepository represents location repository.
type LocationRepository interface {
	// SetLocation sets user's location by given user ID.
	// Returns updated location entity.
	SetLocation(ctx context.Context, arg SetLocationArg) (domain.Location, error)
}
