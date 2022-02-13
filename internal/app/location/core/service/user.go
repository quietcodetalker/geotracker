package service

import (
	"context"
	"fmt"
	log2 "log"
	"time"

	"gitlab.com/spacewalker/geotracker/internal/app/location/core/domain"
	"gitlab.com/spacewalker/geotracker/internal/app/location/core/port"
	"gitlab.com/spacewalker/geotracker/internal/pkg/errpack"
	"gitlab.com/spacewalker/geotracker/internal/pkg/geo"
	"gitlab.com/spacewalker/geotracker/internal/pkg/log"
	"gitlab.com/spacewalker/geotracker/internal/pkg/util"
	"gitlab.com/spacewalker/geotracker/internal/pkg/util/pagination"
)

type userService struct {
	repo          port.UserRepository
	historyClient port.HistoryClient
	logger        log.Logger
}

// NewUserService creates instance of UserService and returns its pointer.
func NewUserService(
	repo port.UserRepository,
	historyClient port.HistoryClient,
	logger log.Logger,
) port.UserService {
	if logger == nil {
		log2.Panic("logger must not be nil")
	}
	if repo == nil {
		logger.Panic("repo must not be nil", nil)
	}
	if historyClient == nil {
		logger.Panic("historyClient must not be nil", nil)
	}

	return &userService{
		repo:          repo,
		historyClient: historyClient,
		logger:        logger,
	}
}

// SetUserLocation sets user's location by given username.
func (s *userService) SetUserLocation(ctx context.Context, req port.UserServiceSetUserLocationRequest) (port.UserServiceSetUserLocationResponse, error) {
	var err error
	defer func() {
		util.LogInternalError(ctx, s.logger, err, req)
	}()

	if err = validate.Struct(req); err != nil {
		// TODO: check different errors?
		return port.UserServiceSetUserLocationResponse{}, fmt.Errorf("%w", errpack.ErrInvalidArgument)
	}

	point := geo.Trunc(geo.Point{req.Longitude, req.Latitude})

	res, err := s.repo.SetUserLocation(ctx, port.UserRepositorySetUserLocationRequest{
		Username: req.Username,
		Point:    point,
	})
	if err != nil {
		return port.UserServiceSetUserLocationResponse{}, err
	}

	if res.PrevLocation.UserID == res.User.ID {
		_, _ = s.historyClient.AddRecord(ctx, port.HistoryClientAddRecordRequest{
			UserID:    res.PrevLocation.UserID,
			A:         res.PrevLocation.Point,
			B:         res.Location.Point,
			Timestamp: time.Now().UTC(),
		})
	}

	return port.UserServiceSetUserLocationResponse{
		Latitude:  res.Location.Point.Latitude(),
		Longitude: res.Location.Point.Longitude(),
	}, nil
}

// ListUsersInRadius finds users by given location and radius.
func (s *userService) ListUsersInRadius(ctx context.Context, req port.UserServiceListUsersInRadiusRequest) (port.UserServiceListUsersInRadiusResponse, error) {
	var err error
	defer func() {
		util.LogInternalError(ctx, s.logger, err, req)
	}()

	if err = validate.Struct(req); err != nil {
		// TODO: check different errors
		return port.UserServiceListUsersInRadiusResponse{}, fmt.Errorf("%w", errpack.ErrInvalidArgument)
	}

	req.Point = geo.Trunc(req.Point)

	var pageToken, pageSize int
	if req.PageToken != "" {
		var err error
		pageToken, pageSize, err = pagination.DecodeCursor(req.PageToken)
		if err != nil {
			return port.UserServiceListUsersInRadiusResponse{}, fmt.Errorf("%w", errpack.ErrInvalidArgument)
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

	if res.Users == nil {
		res.Users = make([]domain.User, 0)
	}

	return port.UserServiceListUsersInRadiusResponse{
		Users:         res.Users,
		NextPageToken: nextPageToken,
	}, nil
}

// GetByUsername finds user by username.
//
// It returns a user and any error encountered.
//
// `ErrInvalidArgument` is returned in case username is empty string.
//
// Any other error occurred in `GetByUsername` is returned.
func (s *userService) GetByUsername(ctx context.Context, username string) (domain.User, error) {
	var err error
	defer func() {
		util.LogInternalError(ctx, s.logger, err, username)
	}()

	if username == "" {
		return domain.User{}, fmt.Errorf("%w", errpack.ErrInvalidArgument)
	}

	user, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}
