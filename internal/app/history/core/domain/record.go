package domain

import (
	"time"

	"gitlab.com/spacewalker/geotracker/internal/pkg/geo"
)

// Record represents history record of users` movements.
type Record struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	A         geo.Point `json:"a"`
	B         geo.Point `json:"b"`
	Timestamp time.Time `json:"timestamp"`
}
