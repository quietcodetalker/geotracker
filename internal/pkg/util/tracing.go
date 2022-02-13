package util

import (
	"context"
	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
	"strconv"
	"time"
)

// TraceIDCtxKey is a type key for trace id.
type TraceIDCtxKey = struct{}

// TraceIDMetadataKey TODO: description
const TraceIDMetadataKey = "trace-id"

// GenerateTraceID returns a new trace id as a string.
//
// It returns uuid as the id. In case uuid generation is failed,
// it returns unix timestamp.
func GenerateTraceID() string {
	uuid, err := uuid.NewUUID()
	if err != nil {
		return strconv.FormatInt(time.Now().Unix(), 10)
	}

	return uuid.String()
}

// AddTraceIDToCtx adds trace id to ctx and returns it.
func AddTraceIDToCtx(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, TraceIDCtxKey{}, id)
}

// GetTraceIDFromCtx retrieves trace id from ctx.
//
// It returns id and ok which is true if token exists and false otherwise.
func GetTraceIDFromCtx(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(TraceIDCtxKey{}).(string)
	return id, ok
}

// AddTraceIDToMetadata TODO: description
func AddTraceIDToMetadata(ctx context.Context, id string) context.Context {
	md := metadata.Pairs(TraceIDMetadataKey, id)
	mdCtx := metadata.NewOutgoingContext(ctx, md)

	return mdCtx
}

// GetTraceIDFromMetadata TODO: description
func GetTraceIDFromMetadata(ctx context.Context) (string, bool) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", false
	}

	value, ok := md[TraceIDMetadataKey]
	if !ok {
		return "", false
	}
	if len(value) == 0 {
		return "", false
	}

	return value[0], true
}
