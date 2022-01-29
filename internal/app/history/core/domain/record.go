package domain

import (
	"gitlab.com/spacewalker/locations/internal/pkg/geo"
	"time"
)

// Record represents history record of users` movements.
type Record struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	A         geo.Point `json:"a"`
	B         geo.Point `json:"b"`
	Timestamp time.Time `json:"timestamp"`
}
