package handler

import (
	"context"
	"gitlab.com/spacewalker/locations/internal/app/location/core/port"
	"gitlab.com/spacewalker/locations/internal/pkg/errpack"
	pb "gitlab.com/spacewalker/locations/pkg/api/proto/v1/location"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// GRPCHandler is a handler that serves request from another internal microservice.
type GRPCHandler struct {
	service port.UserService

	pb.UnimplementedLocationInternalServer
}

// NewGRPCHandler creates new instance of GRPCHandler and returns its pointer.
func NewGRPCHandler(service port.UserService) *GRPCHandler {
	return &GRPCHandler{
		service: service,
	}
}

// GetUserByUsername finds user by username.
func (h *GRPCHandler) GetUserByUsername(ctx context.Context, request *pb.GetUserByUsernameRequest) (*pb.User, error) {
	user, err := h.service.GetByUsername(ctx, request.Username)
	if err != nil {
		return nil, errpack.ErrToGRPC(err)
	}

	return &pb.User{
		Id:        int32(user.ID),
		Username:  user.Username,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}, errpack.ErrToGRPC(nil)
}
