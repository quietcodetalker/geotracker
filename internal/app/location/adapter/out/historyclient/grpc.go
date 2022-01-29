package historyclient

import (
	"context"
	"fmt"
	"gitlab.com/spacewalker/locations/internal/app/location/core/port"
	"gitlab.com/spacewalker/locations/internal/pkg/errpack"
	"gitlab.com/spacewalker/locations/internal/pkg/geo"
	pb "gitlab.com/spacewalker/locations/pkg/api/proto/v1/history"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCClient struct {
	addr string
}

// NewGRPCClient TODO: add description
func NewGRPCClient(addr string) port.HistoryClient {
	return &GRPCClient{
		addr: addr,
	}
}

// AddRecord TODO: add description
func (c GRPCClient) AddRecord(ctx context.Context, req port.HistoryClientAddRecordRequest) (port.HistoryClientAddRecordResponse, error) {
	var opts []grpc.DialOption

	opts = append(opts, grpc.WithInsecure())

	conn, err := grpc.Dial(c.addr, opts...)
	if err != nil {
		return port.HistoryClientAddRecordResponse{}, fmt.Errorf("%w", errpack.ErrInternalError)
	}
	defer conn.Close()

	client := pb.NewHistoryClient(conn)

	res, err := client.AddRecord(ctx, &pb.AddRecordRequest{
		UserId: int32(req.UserID),
		A: &pb.Point{
			Longitude: req.A.Longitude(),
			Latitude:  req.A.Latitude(),
		},
		B: &pb.Point{
			Longitude: req.B.Longitude(),
			Latitude:  req.B.Latitude(),
		},
	})
	if err != nil {
		st, ok := status.FromError(err)
		if !ok {
			return port.HistoryClientAddRecordResponse{}, fmt.Errorf("%w", errpack.ErrInternalError)
		}
		switch st.Code() {
		case codes.InvalidArgument:
			return port.HistoryClientAddRecordResponse{}, fmt.Errorf("%w", errpack.ErrInvalidArgument)
		default:
			return port.HistoryClientAddRecordResponse{}, fmt.Errorf("%w", errpack.ErrInternalError)
		}
	}

	return port.HistoryClientAddRecordResponse{
		UserID:    int(res.UserId),
		A:         geo.Point{res.A.Longitude, res.A.Latitude},
		B:         geo.Point{res.B.Longitude, res.B.Latitude},
		CreatedAt: res.CreatedAt.AsTime(),
		UpdatedAt: res.UpdatedAt.AsTime(),
	}, nil
}

func (c GRPCClient) GetDistance(ctx context.Context, req port.HistoryClientGetDistanceRequest) (port.HistoryClientGetDistanceResponse, error) {
	panic("implement me")
}
