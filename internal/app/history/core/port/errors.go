package port

import (
	"errors"
)

var (
	// ErrNotFound error means that requested entity is not found.
	ErrNotFound = errors.New("not found")
)

// TODO add description
type InvalidLocationErrorViolation struct {
	Subject string
	Value   float64
}

// TODO add description
type InvalidLocationError struct {
	Violations []InvalidLocationErrorViolation
}

// TODO add description
func (e *InvalidLocationError) Error() string {
	return "invalid location"
}

// TODO add description
type InvalidArgumentError struct{}

// TODO add description
func (e *InvalidArgumentError) Error() string {
	return "invalid argument"
}
