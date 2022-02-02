package middleware

import (
	"context"
	"fmt"
	"gitlab.com/spacewalker/locations/internal/pkg/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"time"
)

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

		m, err := handler(ctx, req)

		st, _ := status.FromError(err)

		logger.Info("grpc request complete", log.Fields{
			"method":   fullMethod,
			"duration": time.Since(start),
			"code":     st.Code(),
			"request":  request,
		})

		return m, err
	}
}
