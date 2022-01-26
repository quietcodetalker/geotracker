package port

import (
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

// InternalError represents internal error.
type InternalError struct {
}

// Error returns text representation of the error.
func (e *InternalError) Error() string {
	return "internal error"
}

// ErrToGRPC converts err to Error from `google.golang.org/grpc/status`.
func ErrToGRPC(err error) error {
	var internalError *InternalError
	var invalidLocationError *InvalidLocationError
	var invalidArgumentError *InvalidArgumentError

	switch {
	case err == nil:
		return status.Error(codes.OK, "")
	case errors.As(err, &internalError):
		return status.Error(codes.Internal, err.Error())
	case errors.As(err, &invalidLocationError):
		return status.Error(codes.FailedPrecondition, err.Error())
	case errors.As(err, &invalidArgumentError):
		return status.Error(codes.FailedPrecondition, err.Error())
	case errors.Is(err, ErrNotFound):
		return status.Error(codes.NotFound, err.Error())
	default:
		return status.Error(codes.Unknown, "")
	}
}

// ErrToHTTP converts err to representation of JSON object.
func ErrToHTTP(err error) map[string]interface{} {
	var internalError *InternalError
	var invalidLocationError *InvalidLocationError
	var invalidArgumentError *InvalidArgumentError

	switch {
	case err == nil:
		return map[string]interface{}{
			"error": map[string]interface{}{
				"code":    200,
				"message": "OK",
				"status":  "OK",
			},
		}
	case errors.As(err, &internalError):
		return map[string]interface{}{
			"error": map[string]interface{}{
				"code":    500,
				"message": "internal error",
				"status":  "INTERNAL",
			},
		}
	case errors.As(err, &invalidLocationError):
		return map[string]interface{}{
			"error": map[string]interface{}{
				"code":    400,
				"message": "bad request",
				"status":  "FAILED_PRECONDITION",
			},
		}
	case errors.As(err, &invalidArgumentError):
		return map[string]interface{}{
			"error": map[string]interface{}{
				"code":    400,
				"message": "bad request",
				"status":  "FAILED_PRECONDITION",
			},
		}
	case errors.Is(err, ErrNotFound):
		return map[string]interface{}{
			"error": map[string]interface{}{
				"code":    404,
				"message": "not found",
				"status":  "NOT_FOUND",
			},
		}
	default:
		return map[string]interface{}{
			"error": map[string]interface{}{
				"code":    500,
				"message": "unknown error",
				"status":  "UNKNOWN",
			},
		}
	}
}
