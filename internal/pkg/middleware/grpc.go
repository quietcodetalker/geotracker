package middleware

import (
	"context"
	"fmt"
	"gitlab.com/spacewalker/locations/internal/pkg/log"
	"gitlab.com/spacewalker/locations/internal/pkg/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"time"
)

// LoggerUnaryServerInterceptor TODO: description
func LoggerUnaryServerInterceptor(logger log.Logger) func(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()
		fullMethod := info.FullMethod
		request := fmt.Sprintf("%v", req)

		traceID, ok := util.GetTraceIDFromMetadata(ctx)
		if !ok {
			traceID = util.GenerateTraceID()
		}
		util.AddTraceIDToCtx(ctx, traceID)

		m, err := handler(ctx, req)

		st, _ := status.FromError(err)

		logger.Info("incoming grpc request complete", log.Fields{
			"method":   fullMethod,
			"duration": time.Since(start),
			"code":     st.Code(),
			"request":  request,
			"trace-id": traceID,
		})

		return m, err
	}
}
