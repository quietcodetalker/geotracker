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
		mdCtx := util.AddTraceIDToCtx(ctx, traceID)

		m, err := handler(mdCtx, req)

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

// LoggerUnaryClientInterceptor TODO: description
func LoggerUnaryClientInterceptor(logger log.Logger) func(
	ctx context.Context,
	method string,
	req interface{},
	reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	return func(
		ctx context.Context,
		method string,
		req interface{},
		reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		mdCtx := ctx

		start := time.Now()

		traceID, ok := util.GetTraceIDFromCtx(ctx)
		if !ok {
			logger.Warn("traceID is nil", log.Fields{
				"method":  method,
				"request": req,
			})
		} else {
			mdCtx = util.AddTraceIDToMetadata(ctx, traceID)
		}

		err := invoker(mdCtx, method, req, reply, cc, opts...)

		st, _ := status.FromError(err)

		logger.Info("outgoing grpc request complete", log.Fields{
			"method":   method,
			"duration": time.Since(start),
			"code":     st.Code(),
			"request":  req,
			"trace-id": traceID,
		})

		return err
	}
}
