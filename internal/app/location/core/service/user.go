package service

import (
	"context"
	"gitlab.com/spacewalker/locations/internal/app/location/core/domain"
	"gitlab.com/spacewalker/locations/internal/app/location/core/port"
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
func (s *userService) SetUserLocation(ctx context.Context, req port.SetUserLocationRequest) (port.SetUserLocationResponse, error) {
	if err := validate.Struct(req); err != nil {
		// TODO: check different errors?
		return port.SetUserLocationResponse{}, &port.InvalidArgumentError{}
	}

	location, err := s.repo.SetUserLocation(ctx, port.SetUserLocationArg{
		Username: req.Username,
		Point:    domain.Point{req.Longitude, req.Latitude},
	})
	if err != nil {
		return port.SetUserLocationResponse{}, err
	}

	return port.SetUserLocationResponse{
		Latitude:  location.Point.Latitude(),
		Longitude: location.Point.Longitude(),
	}, nil

	// TODO: add record to location history via history microservice
}
