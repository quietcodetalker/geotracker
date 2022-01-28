package errpack

import (
	"errors"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// PreconditionFailureViolation describes a single precondition failure.
type PreconditionFailureViolation struct {
	Type        string `json:"type"`
	Subject     string `json:"subject"`
	Description string `json:"description"`
}

// PreconditionFailureError describes failed preconditions.
type PreconditionFailureError struct {
	Violations []*PreconditionFailureViolation `json:"violations"`
	Err        error
}

// Error returns error's text representation.
func (e *PreconditionFailureError) Error() string {
	return "precondition failure"
}

// BadRequestViolation describes a single bad request field.
type BadRequestViolation struct {
	Field       string `json:"field"`
	Description string `json:"description"`
}

// BadRequestError describes request violations.
type BadRequestError struct {
	Violations []*BadRequestViolation `json:"violations"`
	Err        error
}

// Error returns error's text representation.
func (e *BadRequestError) Error() string {
	return "bad request"
}

// ResourceInfo describes the resource is being accessed.
type ResourceInfo struct {
	ResourceType string `json:"resource_type"`
	ResourceName string `json:"resource_name"`
	Owner        string `json:"owner"`
	Description  string `json:"description"`
}

// NotFoundError means that resource is being accessed isn't found.
type NotFoundError struct {
	ResourceInfo
	Err error
}

// Error returns error's text representation.
func (e *NotFoundError) Error() string {
	return e.Description
}

// AlreadyExistsError means that resource is being created already exists.
type AlreadyExistsError struct {
	ResourceInfo
	Err error
}

// Error returns error's text representation.
func (e *AlreadyExistsError) Error() string {
	return e.Description
}

// DebugInfo describes debug info.
type DebugInfo struct {
	StackEntries []string `json:"stack_entries"`
	Detail       string   `json:"detail"`
}

// InternalError is an internal error.
type InternalError struct {
	DebugInfo
	Err error
}

// Error returns error's text representation.
func (e *InternalError) Error() string {
	return e.Detail
}

// UnknownError is an unknown error.
type UnknownError struct {
	DebugInfo
	Err error
}

// Error returns error's text representation.
func (e *UnknownError) Error() string {
	return e.Detail
}

// GRPCStatusFromErr maps errors to grpc status.
func GRPCStatusFromErr(err error) error {
	var preconditionFailureError *PreconditionFailureError
	var badRequestError *BadRequestError
	var notFoundError *NotFoundError
	var alreadyExistError *AlreadyExistsError
	var internalError *InternalError
	var unknownError *UnknownError

	var st *status.Status
	var stErr error

	switch {
	case errors.As(err, &preconditionFailureError):
		violations := make([]*errdetails.PreconditionFailure_Violation, 0, len(preconditionFailureError.Violations))
		for _, v := range preconditionFailureError.Violations {
			violations = append(violations, &errdetails.PreconditionFailure_Violation{
				Type:        v.Type,
				Subject:     v.Subject,
				Description: v.Description,
			})
		}
		st, stErr = status.New(codes.FailedPrecondition, preconditionFailureError.Error()).
			WithDetails(&errdetails.PreconditionFailure{
				Violations: violations,
			})
	case errors.As(err, &badRequestError):
		violations := make([]*errdetails.BadRequest_FieldViolation, 0, len(badRequestError.Violations))
		for _, v := range badRequestError.Violations {
			violations = append(violations, &errdetails.BadRequest_FieldViolation{
				Field:       v.Field,
				Description: v.Description,
			})
		}
		st, stErr = status.New(codes.InvalidArgument, badRequestError.Error()).
			WithDetails(&errdetails.BadRequest{
				FieldViolations: violations,
			})
	case errors.As(err, &notFoundError):
		st, stErr = status.New(codes.NotFound, notFoundError.Error()).
			WithDetails(&errdetails.ResourceInfo{
				ResourceType: notFoundError.ResourceType,
				ResourceName: notFoundError.ResourceName,
				Owner:        notFoundError.Owner,
				Description:  notFoundError.Description,
			})
	case errors.As(err, &alreadyExistError):
		st, stErr = status.New(codes.AlreadyExists, alreadyExistError.Error()).
			WithDetails(&errdetails.ResourceInfo{
				ResourceType: alreadyExistError.ResourceType,
				ResourceName: alreadyExistError.ResourceName,
				Owner:        alreadyExistError.Owner,
				Description:  alreadyExistError.Description,
			})
	case errors.As(err, &internalError):
		st, stErr = status.New(codes.Internal, internalError.Error()).
			WithDetails(&errdetails.DebugInfo{
				StackEntries: internalError.StackEntries,
				Detail:       internalError.Detail,
			})
	case errors.As(err, &unknownError):
		st, stErr = status.New(codes.Unknown, unknownError.Error()).
			WithDetails(&errdetails.DebugInfo{
				StackEntries: unknownError.StackEntries,
				Detail:       unknownError.Detail,
			})
	default:
		st, stErr = status.New(codes.Unknown, unknownError.Error()).
			WithDetails(&errdetails.DebugInfo{
				Detail: "unknown error",
			})
	}

	if stErr != nil {
		return status.Error(codes.Unknown, "unknown error")
	}
	return st.Err()
}
