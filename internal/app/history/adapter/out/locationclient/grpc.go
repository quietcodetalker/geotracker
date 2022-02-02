package locationclient

import (
	"context"
	"fmt"
	"gitlab.com/spacewalker/locations/internal/pkg/errpack"
	"gitlab.com/spacewalker/locations/internal/pkg/log"
	"gitlab.com/spacewalker/locations/internal/pkg/middleware"
	pb "gitlab.com/spacewalker/locations/pkg/api/proto/v1/location"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	log2 "log"
)

type GRPCClient struct {
	addr   string
	logger log.Logger
}

// NewGRPCClient TODO: add description
func NewGRPCClient(addr string, logger log.Logger) *GRPCClient {
	if logger == nil {
		log2.Panic("logger must not be nil")
	}

	return &GRPCClient{
		addr:   addr,
		logger: logger,
	}
}

// GetUserIDByUsername TODO: add description
func (c *GRPCClient) GetUserIDByUsername(ctx context.Context, username string) (int, error) {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(middleware.LoggerUnaryClientInterceptor(c.logger)),
	}

	grpc.WithInsecure()

	conn, err := grpc.Dial(c.addr, opts...)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", errpack.ErrInternalError, err)
	}
	defer conn.Close()

	client := pb.NewLocationInternalClient(conn)

	user, err := client.GetUserByUsername(ctx, &pb.GetUserByUsernameRequest{Username: username})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return 0, fmt.Errorf("%w: %v", errpack.ErrInternalError, err)
		}
		switch st.Code() {
		case codes.NotFound:
			return 0, fmt.Errorf("%w", errpack.ErrNotFound)
		default:
			return 0, fmt.Errorf("%w: %v", errpack.ErrInternalError, err)
		}
	}

	return int(user.Id), nil
}
