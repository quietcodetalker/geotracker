package util

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"gitlab.com/spacewalker/locations/internal/pkg/errpack"
)

// TraceIDCtxKey is a type key for trace id.
type TraceIDCtxKey = struct{}

// GenerateTraceID returns uuid as a string.
//
// Returns a generated id and any error encountered.
//
// Returned error wraps `ErrInternalError` in case of any error occurred.
func GenerateTraceID() (string, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return "", fmt.Errorf("%w: %v", errpack.ErrInternalError, err)
	}

	return id.String(), nil
}

// AddTraceID adds trace id to ctx and returns it.
func AddTraceID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, TraceIDCtxKey{}, id)
}

// GetTraceID retrieves trace id from ctx.
// If ctx does cot contain trace id, it generates new traced id and adds the id to ctx.
//
// It returns provided or updated context, retrieved or generated trace id and error.
//
// In case of any error while generating new trace id with GenerateTraceID,
// provided context without changes, empty string and the error are returned.
func GetTraceID(ctx context.Context) (context.Context, string, error) {
	id, ok := ctx.Value(TraceIDCtxKey{}).(string)
	if !ok {
		id, err := GenerateTraceID()
		if err != nil {
			return ctx, "", err
		}
		return AddTraceID(ctx, id), id, nil
	}

	return ctx, id, nil
}
