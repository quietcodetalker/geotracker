package service

import (
	"context"
	"gitlab.com/spacewalker/locations/internal/app/history/core/port"
	"gitlab.com/spacewalker/locations/internal/pkg/geo"
	"time"
)

type historyService struct {
	repo           port.HistoryRepository
	locationClient port.LocationClient
}

// NewHistoryService creates new instance of history service and returns its pointer.
func NewHistoryService(repo port.HistoryRepository, locationClient port.LocationClient) port.HistoryService {
	return &historyService{
		repo:           repo,
		locationClient: locationClient,
	}
}

// AddRecord adds a history record via repository and returns it (except id).
func (s *historyService) AddRecord(ctx context.Context, req port.HistoryServiceAddRecordRequest) (port.HistoryServiceAddRecordResponse, error) {
	if err := validate.Struct(req); err != nil {
		// TODO: handle different errors
		return port.HistoryServiceAddRecordResponse{}, &port.InvalidArgumentError{}
	}

	record, err := s.repo.AddRecord(ctx, port.HistoryRepositoryAddRecordRequest{
		UserID: req.UserID,
		A:      geo.Trunc(req.A),
		B:      geo.Trunc(req.B),
	})
	if err != nil {
		return port.HistoryServiceAddRecordResponse{}, err
	}

	return port.HistoryServiceAddRecordResponse{
		UserID:    record.UserID,
		A:         record.A,
		B:         record.B,
		CreatedAt: record.CreatedAt,
		UpdatedAt: record.UpdatedAt,
	}, nil
}

// GetDistance calculates distance that particular user got through in given time period.
func (s *historyService) GetDistance(ctx context.Context, req port.HistoryServiceGetDistanceRequest) (port.HistoryServiceGetDistanceResponse, error) {
	if err := validate.Struct(req); err != nil {
		return port.HistoryServiceGetDistanceResponse{}, &port.InvalidArgumentError{}
	}

	distance, err := s.repo.GetDistance(ctx, port.HistoryRepositoryGetDistanceRequest(req))
	if err != nil {
		return port.HistoryServiceGetDistanceResponse{}, err
	}

	return port.HistoryServiceGetDistanceResponse{Distance: distance}, nil
}

// GetDistanceByUsername calculates distance that particular user got through in given time period.
func (s *historyService) GetDistanceByUsername(ctx context.Context, req port.HistoryServiceGetDistanceByUsernameRequest) (port.HistoryServiceGetDistanceByUsernameResponse, error) {
	if err := validate.Struct(req); err != nil {
		return port.HistoryServiceGetDistanceByUsernameResponse{}, &port.InvalidArgumentError{}
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
