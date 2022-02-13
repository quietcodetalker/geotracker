package domain

import (
	"time"

	"gitlab.com/spacewalker/geotracker/internal/pkg/geo"
)

// Location represents user's geographic position.
type Location struct {
	UserID    int       `json:"user_id"`
	Point     geo.Point `json:"point"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
