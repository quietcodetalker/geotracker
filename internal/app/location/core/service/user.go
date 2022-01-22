package service

import (
	"context"
	"gitlab.com/spacewalker/locations/internal/app/location/core/port"
	"gitlab.com/spacewalker/locations/internal/pkg/geo"
	"gitlab.com/spacewalker/locations/internal/pkg/util/pagination"
)

type userService struct {
	repo port.UserRepository
}

// NewUserService creates instance of UserService and returns its pointer.
func NewUserService(repo port.UserRepository) port.UserService {
	return &userService{
		repo: repo,
	}
}

// SetUserLocation sets user's location by given username.
func (s *userService) SetUserLocation(ctx context.Context, req port.UserServiceSetUserLocationRequest) (port.UserServiceSetUserLocationResponse, error) {
	if err := validate.Struct(req); err != nil {
		// TODO: check different errors?
		return port.UserServiceSetUserLocationResponse{}, &port.InvalidArgumentError{}
	}

	location, err := s.repo.SetUserLocation(ctx, port.UserRepositorySetUserLocationRequest{
		Username: req.Username,
		Point:    geo.Point{req.Longitude, req.Latitude},
	})
	if err != nil {
		return port.UserServiceSetUserLocationResponse{}, err
	}

	return port.UserServiceSetUserLocationResponse{
		Latitude:  location.Point.Latitude(),
		Longitude: location.Point.Longitude(),
	}, nil

	// TODO: add record to location history via history microservice
}

// ListUsersInRadius finds users around given geographic point in given radius.
func (s *userService) ListUsersInRadius(ctx context.Context, req port.UserServiceListUsersInRadiusRequest) (port.UserServiceListUsersInRadiusResponse, error) {
	if err := validate.Struct(req); err != nil {
		// TODO: check different errors
		return port.UserServiceListUsersInRadiusResponse{}, &port.InvalidArgumentError{}
	}

	var pageToken, pageSize int
	if req.PageToken != "" {
		var err error
		pageToken, pageSize, err = pagination.DecodeCursor(req.PageToken)
		if err != nil {
			return port.UserServiceListUsersInRadiusResponse{}, &port.InvalidArgumentError{}
		}
	} else {
		pageSize = req.PageSize
	}

	res, err := s.repo.ListUsersInRadius(ctx, port.UserRepositoryListUsersInRadiusRequest{
		Point:     req.Point,
		Radius:    req.Radius,
		PageToken: pageToken,
		PageSize:  pageSize,
	})
	if err != nil {
		return port.UserServiceListUsersInRadiusResponse{}, err
	}

	nextPageToken := ""
	if res.NextPageToken > 0 {
		nextPageToken = pagination.EncodeCursor(res.NextPageToken, pageSize)
	}

	return port.UserServiceListUsersInRadiusResponse{
		Users:         res.Users,
		NextPageToken: nextPageToken,
	}, nil
}
