package handler

import (
	"context"
	"fmt"

	"gitlab.com/spacewalker/geotracker/internal/app/history/core/port"
	"gitlab.com/spacewalker/geotracker/internal/pkg/errpack"
	"gitlab.com/spacewalker/geotracker/internal/pkg/geo"
	pb "gitlab.com/spacewalker/geotracker/pkg/api/proto/v1/history"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
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

// AddRecord TODO: description
func (h *GRPCHandler) AddRecord(ctx context.Context, req *pb.AddRecordRequest) (*pb.AddRecordResponse, error) {
	if req.A == nil || req.B == nil || req.Timestamp == nil {
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
		Timestamp: req.Timestamp.AsTime(),
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
		Timestamp: timestamppb.New(res.Timestamp),
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
