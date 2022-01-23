package handler

import (
	"context"
	"errors"
	"gitlab.com/spacewalker/locations/internal/app/location/core/port"
	"gitlab.com/spacewalker/locations/internal/pkg/geo"
	pb "gitlab.com/spacewalker/locations/pkg/api/proto/v1/location"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"net"
)

type GRPCHandler struct {
	server  *grpc.Server
	service port.UserService

	pb.UnimplementedLocationServer
}

// NewGRPCHandler creates an instance of GRPCHandler and returns its pointer.
func NewGRPCHandler(service port.UserService) *GRPCHandler {
	return &GRPCHandler{
		service: service,
	}
}

// Start starts grpc server.
func (h *GRPCHandler) Start(port string) error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}

	h.server = grpc.NewServer()

	pb.RegisterLocationServer(h.server, h)

	if err := h.server.Serve(lis); err != nil {
		if err == grpc.ErrServerStopped {
			return nil
		}

		return err
	}

	return nil
}

// Stop stops grpc server.
func (h *GRPCHandler) Stop() {
	h.server.Stop()
}

// SetUserLocation TODO: add description
func (h *GRPCHandler) SetUserLocation(
	ctx context.Context,
	request *pb.SetUserLocationRequest,
) (*pb.SetUserLocationResponse, error) {
	response, err := h.service.SetUserLocation(ctx, port.UserServiceSetUserLocationRequest{
		Username:  request.GetUsername(),
		Latitude:  request.GetLatitude(),
		Longitude: request.GetLongitude(),
	})
	if err != nil {
		return nil, grpcErr(err)
	}

	return &pb.SetUserLocationResponse{
		Longitude: response.Longitude,
		Latitude:  response.Latitude,
	}, status.Error(codes.OK, "")
}

// ListUsersInRadius TODO: add description
func (h *GRPCHandler) ListUsersInRadius(
	ctx context.Context,
	req *pb.ListUsersInRadiusRequest,
) (*pb.ListUsersInRadiusResponse, error) {
	if len(req.Point) != 2 {
		return &pb.ListUsersInRadiusResponse{}, status.Error(codes.FailedPrecondition, "invalid point")
	}
	res, err := h.service.ListUsersInRadius(
		ctx,
		port.UserServiceListUsersInRadiusRequest{
			Point:     geo.Point{req.Point[0], req.Point[1]},
			Radius:    req.Radius,
			PageToken: req.PageToken,
			PageSize:  int(req.PageSize),
		},
	)
	if err != nil {
		return nil, grpcErr(err)
	}

	users := make([]*pb.User, 0, len(res.Users))
	for _, user := range res.Users {
		users = append(users, &pb.User{
			Id:        int32(user.ID),
			Username:  user.Username,
			CreatedAt: timestamppb.New(user.CreatedAt),
			UpdatedAt: timestamppb.New(user.UpdatedAt),
		})
	}

	return &pb.ListUsersInRadiusResponse{
		Users:         users,
		NextPageToken: res.NextPageToken,
	}, status.Error(codes.OK, "")
}

func grpcErr(err error) error {
	var invalidArgumentError *port.InvalidArgumentError

	if err != nil {
		switch {
		case errors.As(err, &invalidArgumentError):
			fallthrough
		case errors.Is(err, port.ErrInvalidUsername):
			return status.Error(codes.FailedPrecondition, err.Error())
		default:
			return status.Error(codes.Internal, "internal error")
		}
	}

	return nil
}
