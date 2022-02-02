package locationclient

import (
	"context"
	"fmt"
	"gitlab.com/spacewalker/locations/internal/pkg/errpack"
	pb "gitlab.com/spacewalker/locations/pkg/api/proto/v1/location"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCClient struct {
	addr string
}

// NewGRPCClient TODO: add description
func NewGRPCClient(addr string) *GRPCClient {
	return &GRPCClient{addr: addr}
}

// GetUserIDByUsername TODO: add description
func (c *GRPCClient) GetUserIDByUsername(ctx context.Context, username string) (int, error) {
	var opts []grpc.DialOption

	opts = append(opts, grpc.WithInsecure())

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
