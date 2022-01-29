package errpack

import (
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ErrToGRPC converts err to Error from `google.golang.org/grpc/status`.
func ErrToGRPC(err error) error {
	switch {
	case err == nil:
		return status.Error(codes.OK, "OK")
	case errors.Is(err, ErrInternalError):
		return status.Error(codes.Internal, err.Error())
	case errors.Is(err, ErrInvalidArgument):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, ErrNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, ErrFailedPrecondition):
		return status.Error(codes.FailedPrecondition, err.Error())
	default:
		return status.Error(codes.Unknown, "unknown error")
	}
}

// ErrToHTTP converts err to representation of JSON object.
func ErrToHTTP(err error) map[string]interface{} {
	switch {
	case err == nil:
		return map[string]interface{}{
			"error": map[string]interface{}{
				"code":    200,
				"message": "OK",
				"status":  "OK",
			},
		}
	case errors.Is(err, ErrInternalError):
		return map[string]interface{}{
			"error": map[string]interface{}{
				"code":    500,
				"message": err.Error(),
				"status":  "INTERNAL",
			},
		}
	case errors.Is(err, ErrFailedPrecondition):
		return map[string]interface{}{
			"error": map[string]interface{}{
				"code":    422,
				"message": err.Error(),
				"status":  "FAILED_PRECONDITION",
			},
		}
	case errors.Is(err, ErrInvalidArgument):
		return map[string]interface{}{
			"error": map[string]interface{}{
				"code":    400,
				"message": err.Error(),
				"status":  "INVALID_ARGUMENT",
			},
		}
	case errors.Is(err, ErrNotFound):
		return map[string]interface{}{
			"error": map[string]interface{}{
				"code":    404,
				"message": err.Error(),
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
