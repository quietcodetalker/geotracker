package locationclient

import (
	"context"
	pb "gitlab.com/spacewalker/locations/pkg/api/proto/v1/location"
	"google.golang.org/grpc"
)

type GRPCClient struct {
	addr string
}

func NewGRPCClient(addr string) *GRPCClient {
	return &GRPCClient{addr: addr}
}

func (c *GRPCClient) GetUserIDByUsername(ctx context.Context, username string) (int, error) {
	var opts []grpc.DialOption

	opts = append(opts, grpc.WithInsecure())

	conn, err := grpc.Dial(c.addr, opts...)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	client := pb.NewLocationInternalClient(conn)

	user, err := client.GetUserByUsername(ctx, &pb.GetUserByUsernameRequest{Username: username})
	if err != nil {
		return 0, err
	}

	return int(user.Id), nil
}
