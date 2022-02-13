package service

import (
	"context"
	"fmt"
	log2 "log"
	"time"

	"gitlab.com/spacewalker/geotracker/internal/app/history/core/domain"
	"gitlab.com/spacewalker/geotracker/internal/app/history/core/port"
	"gitlab.com/spacewalker/geotracker/internal/pkg/errpack"
	"gitlab.com/spacewalker/geotracker/internal/pkg/geo"
	"gitlab.com/spacewalker/geotracker/internal/pkg/log"
	"gitlab.com/spacewalker/geotracker/internal/pkg/util"
)

type historyService struct {
	repo           port.HistoryRepository
	locationClient port.LocationClient
	logger         log.Logger
}

// NewHistoryService creates new instance of history service and returns its pointer.
func NewHistoryService(
	repo port.HistoryRepository,
	locationClient port.LocationClient,
	logger log.Logger,
) port.HistoryService {
	if logger == nil {
		log2.Panic("logger must not be nil")
	}
	if repo == nil {
		logger.Panic("repo must not be nil", nil)
	}
	if locationClient == nil {
		logger.Panic("locationClient must not be nil", nil)
	}

	return &historyService{
		repo:           repo,
		locationClient: locationClient,
		logger:         logger,
	}
}

// AddRecord adds a history record.
//
// It returns an added record and any error occurred.
//
// `ErrInvalidArgument` is returned in case of `req` validation failure.
//
// If a call to `AddRecord` repository method fails, any returned error is propagated.
func (s *historyService) AddRecord(ctx context.Context, req port.HistoryServiceAddRecordRequest) (domain.Record, error) {
	var err error
	defer func() {
		util.LogInternalError(ctx, s.logger, err, req)
	}()

	if err = validate.Struct(req); err != nil {
		// TODO: handle different errors
		return domain.Record{}, fmt.Errorf("%w", errpack.ErrInvalidArgument)
	}

	record, err := s.repo.AddRecord(ctx, port.HistoryRepositoryAddRecordRequest{
		UserID:    req.UserID,
		A:         geo.Trunc(req.A),
		B:         geo.Trunc(req.B),
		Timestamp: req.Timestamp,
	})
	if err != nil {
		return domain.Record{}, err
	}

	return record, nil
}

// GetDistance calculates distance that particular user got through in given time period.
func (s *historyService) GetDistance(ctx context.Context, req port.HistoryServiceGetDistanceRequest) (port.HistoryServiceGetDistanceResponse, error) {
	var err error
	defer func() {
		util.LogInternalError(ctx, s.logger, err, req)
	}()

	if err = validate.Struct(req); err != nil {
		return port.HistoryServiceGetDistanceResponse{}, fmt.Errorf("%w", errpack.ErrInvalidArgument)
	}

	distance, err := s.repo.GetDistance(ctx, port.HistoryRepositoryGetDistanceRequest(req))
	if err != nil {
		return port.HistoryServiceGetDistanceResponse{}, err
	}

	return port.HistoryServiceGetDistanceResponse{Distance: distance}, nil
}

// GetDistanceByUsername calculates distance that particular user got through in given time period.
func (s *historyService) GetDistanceByUsername(ctx context.Context, req port.HistoryServiceGetDistanceByUsernameRequest) (port.HistoryServiceGetDistanceByUsernameResponse, error) {
	var err error
	defer func() {
		util.LogInternalError(ctx, s.logger, err, req)
	}()

	if err = validate.Struct(req); err != nil {
		return port.HistoryServiceGetDistanceByUsernameResponse{}, fmt.Errorf("%w", errpack.ErrInvalidArgument)
	}
	switch {
	case req.To == nil && req.From == nil:
		to := time.Now()
		from := to.Add(-24 * time.Hour)
		req.To = &to
		req.From = &from
	case req.To == nil:
		to := req.From.Add(24 * time.Hour)
		req.To = &to
	case req.From == nil:
		from := req.To.Add(-24 * time.Hour)
		req.From = &from
	}

	userID, err := s.locationClient.GetUserIDByUsername(ctx, req.Username)
	if err != nil {
		return port.HistoryServiceGetDistanceByUsernameResponse{}, err
	}

	distance, err := s.repo.GetDistance(ctx, port.HistoryRepositoryGetDistanceRequest{
		UserID: userID,
		To:     *req.To,
		From:   *req.From,
	})
	if err != nil {
		return port.HistoryServiceGetDistanceByUsernameResponse{}, err
	}

	return port.HistoryServiceGetDistanceByUsernameResponse{
		Distance: distance,
	}, nil
}
