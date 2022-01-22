package domain

import (
	"gitlab.com/spacewalker/locations/internal/pkg/geo"
	"time"
)

// Location represents user's geographic position.
type Location struct {
	UserID    int       `json:"user_id"`
	Point     geo.Point `json:"point"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
