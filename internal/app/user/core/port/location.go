package port

import "context"

// SetLocationArg is a param object of user repository SetLocation method.
type SetLocationArg struct {
	UserID    int     `json:"user_id"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type LocationRepo interface {
	// SetLocation sets user's location by given user ID.
	SetLocation(ctx context.Context, arg SetLocationArg) error
}
