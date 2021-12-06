package domain

import (
	"time"
)

// Point represents coordinates [longitude, latitude] of the geographic position.
type Point [2]float64

// Longitude returns longitude of the pont.
func (p *Point) Longitude() float64 {
	return p[0]
}

// Latitude returns latitude of the pont.
func (p *Point) Latitude() float64 {
	return p[1]
}

// Location represents user's geographic position.
type Location struct {
	UserID    int       `json:"user_id"`
	Point     Point     `json:"point"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
