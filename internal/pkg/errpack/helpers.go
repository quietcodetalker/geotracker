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
func ErrToHTTP(err error) (int, map[string]interface{}) {
	switch {
	case err == nil:
		return 200, map[string]interface{}{
			"error": map[string]interface{}{
				"code":    200,
				"message": "OK",
				"status":  "OK",
			},
		}
	case errors.Is(err, ErrInternalError):
		return 500, map[string]interface{}{
			"error": map[string]interface{}{
				"code":    500,
				"message": "internal error",
				"status":  "INTERNAL",
			},
		}
	case errors.Is(err, ErrFailedPrecondition):
		return 422, map[string]interface{}{
			"error": map[string]interface{}{
				"code":    422,
				"message": err.Error(),
				"status":  "FAILED_PRECONDITION",
			},
		}
	case errors.Is(err, ErrInvalidArgument):
		return 400, map[string]interface{}{
			"error": map[string]interface{}{
				"code":    400,
				"message": err.Error(),
				"status":  "INVALID_ARGUMENT",
			},
		}
	case errors.Is(err, ErrNotFound):
		return 404, map[string]interface{}{
			"error": map[string]interface{}{
				"code":    404,
				"message": err.Error(),
				"status":  "NOT_FOUND",
			},
		}
	default:
		return 500, map[string]interface{}{
			"error": map[string]interface{}{
				"code":    500,
				"message": "unknown error",
				"status":  "UNKNOWN",
			},
		}
	}
}
