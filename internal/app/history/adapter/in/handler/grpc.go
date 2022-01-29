package handler

import (
	"context"
	"fmt"
	"gitlab.com/spacewalker/locations/internal/app/history/core/port"
	"gitlab.com/spacewalker/locations/internal/pkg/errpack"
	"gitlab.com/spacewalker/locations/internal/pkg/geo"
	pb "gitlab.com/spacewalker/locations/pkg/api/proto/v1/history"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GRPCHandler represents history handle that handle grpc requests.
type GRPCHandler struct {
	service port.HistoryService

	pb.UnimplementedHistoryServer
}

// NewGRPCHandler creates an instance of history grpc handler and returns its pointer.
func NewGRPCHandler(service port.HistoryService) *GRPCHandler {
	return &GRPCHandler{
		service: service,
	}
}

func (h *GRPCHandler) AddRecord(ctx context.Context, req *pb.AddRecordRequest) (*pb.AddRecordResponse, error) {
	if req.A == nil || req.B == nil {
		// TODO: specify error
		return nil, errpack.ErrToGRPC(fmt.Errorf("%w", errpack.ErrInvalidArgument))
	}

	res, err := h.service.AddRecord(ctx, port.HistoryServiceAddRecordRequest{
		UserID: int(req.UserId),
		A: geo.Point{
			req.A.Longitude,
			req.A.Latitude,
		},
		B: geo.Point{
			req.B.Longitude,
			req.B.Latitude,
		},
	})

	if err != nil {
		return nil, errpack.ErrToGRPC(err)
	}

	return &pb.AddRecordResponse{
		UserId: int32(res.UserID),
		A: &pb.Point{
			Longitude: res.A.Longitude(),
			Latitude:  res.A.Latitude(),
		},
		B: &pb.Point{
			Longitude: res.B.Longitude(),
			Latitude:  res.B.Latitude(),
		},
		CreatedAt: nil,
		UpdatedAt: nil,
	}, status.Error(codes.OK, "")
}

func (h *GRPCHandler) GetDistance(ctx context.Context, req *pb.GetDistanceRequest) (*pb.GetDistanceResponse, error) {
	if req.From == nil || req.To == nil {
		// TODO: specify error
		return nil, errpack.ErrToGRPC(fmt.Errorf("%w", errpack.ErrInvalidArgument))
	}

	res, err := h.service.GetDistance(ctx, port.HistoryServiceGetDistanceRequest{
		UserID: int(req.UserId),
		From:   req.From.AsTime(),
		To:     req.To.AsTime(),
	})
	if err != nil {
		return nil, errpack.ErrToGRPC(err)
	}

	return &pb.GetDistanceResponse{Distance: res.Distance}, status.Error(codes.OK, "")
}
