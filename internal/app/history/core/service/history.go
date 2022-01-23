package service

import (
	"context"
	"fmt"
	"gitlab.com/spacewalker/locations/internal/app/history/core/port"
	"gitlab.com/spacewalker/locations/internal/pkg/geo"
)

type historyService struct {
	repo port.HistoryRepository
}

// NewHistoryService creates new instance of history service and returns its pointer.
func NewHistoryService(repo port.HistoryRepository) port.HistoryService {
	return &historyService{
		repo: repo,
	}
}

// AddRecord adds a history record via repository and returns it (except id).
func (s historyService) AddRecord(ctx context.Context, req port.HistoryServiceAddRecordRequest) (port.HistoryServiceAddRecordResponse, error) {
	if err := validate.Struct(req); err != nil {
		// TODO: handle different errors
		fmt.Println(err)
		return port.HistoryServiceAddRecordResponse{}, &port.InvalidArgumentError{}
	}

	record, err := s.repo.AddRecord(ctx, port.HistoryRepositoryAddRecordRequest{
		UserID: req.UserID,
		A:      geo.Trunc(req.A),
		B:      geo.Trunc(req.B),
	})
	if err != nil {
		fmt.Println(err)
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
func (s historyService) GetDistance(ctx context.Context, req port.HistoryServiceGetDistanceRequest) (port.HistoryServiceGetDistanceResponse, error) {
	if err := validate.Struct(req); err != nil {
		// TODO: handle different errors
		fmt.Println(err)
		return port.HistoryServiceGetDistanceResponse{}, &port.InvalidArgumentError{}
	}

	distance, err := s.repo.GetDistance(ctx, port.HistoryRepositoryGetDistanceRequest(req))
	if err != nil {
		return port.HistoryServiceGetDistanceResponse{}, err
	}

	return port.HistoryServiceGetDistanceResponse{Distance: distance}, nil
}
