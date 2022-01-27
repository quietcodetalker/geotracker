package handler

import (
	"context"
	"gitlab.com/spacewalker/locations/internal/app/location/core/port"
	pb "gitlab.com/spacewalker/locations/pkg/api/proto/v1/location"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// GRPCInternalHandler is a handler that serves request from another internal microservice.
type GRPCInternalHandler struct {
	service port.UserService

	pb.UnimplementedLocationInternalServer
}

// NewGRPCInternalHandler creates new instance of GRPCInternalHandler and returns its pointer.
func NewGRPCInternalHandler(service port.UserService) *GRPCInternalHandler {
	return &GRPCInternalHandler{
		service: service,
	}
}

// GetUserByUsername finds user by username.
func (h *GRPCInternalHandler) GetUserByUsername(ctx context.Context, request *pb.GetUserByUsernameRequest) (*pb.User, error) {
	user, err := h.service.GetByUsername(ctx, request.Username)
	if err != nil {
		if err == port.ErrNotFound {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, "")
	}

	return &pb.User{
		Id:        int32(user.ID),
		Username:  user.Username,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}, status.Error(codes.OK, "")
}
