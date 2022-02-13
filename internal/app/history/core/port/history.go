//go:generate mockgen -destination=mock/mock_user.go -package=mock . HistoryRepository,HistoryService

package port

import (
	"context"
	"time"

	"gitlab.com/spacewalker/geotracker/internal/app/history/core/domain"
	"gitlab.com/spacewalker/geotracker/internal/pkg/geo"
)

// HistoryServiceAddRecordRequest represents request object of HistoryService AddRecord method.
type HistoryServiceAddRecordRequest struct {
	UserID    int       `json:"user_id" validate:"required,gt=0"`
	A         geo.Point `json:"a" validate:"validgeopoint"`
	B         geo.Point `json:"b" validate:"validgeopoint"`
	Timestamp time.Time `json:"timestamp"`
}

// HistoryServiceGetDistanceRequest represents request object of HistoryService GetDistance method.
type HistoryServiceGetDistanceRequest struct {
	UserID int       `json:"user_id" validate:"required,gt=0"`
	From   time.Time `json:"from"`
	To     time.Time `json:"to"`
}

// HistoryServiceGetDistanceResponse represents response object of HistoryService GetDistance method.
type HistoryServiceGetDistanceResponse struct {
	Distance float64 `json:"distance"`
}

// HistoryServiceGetDistanceByUsernameRequest represents request object of HistoryRepository GetDistanceByUsername method.
type HistoryServiceGetDistanceByUsernameRequest struct {
	Username string     `json:"username" validate:"required"`
	From     *time.Time `json:"from"`
	To       *time.Time `json:"to"`
}

// HistoryServiceGetDistanceByUsernameResponse represents response object of HistoryService GetDistanceByUsername method.
type HistoryServiceGetDistanceByUsernameResponse struct {
	Distance float64 `json:"distance"`
}

// HistoryService represents history service.
type HistoryService interface {
	AddRecord(ctx context.Context, req HistoryServiceAddRecordRequest) (domain.Record, error)
	GetDistanceByUsername(ctx context.Context, req HistoryServiceGetDistanceByUsernameRequest) (HistoryServiceGetDistanceByUsernameResponse, error)
	GetDistance(ctx context.Context, req HistoryServiceGetDistanceRequest) (HistoryServiceGetDistanceResponse, error)
}

// HistoryRepositoryAddRecordRequest represents request object of HistoryRepository AddRecord method.
type HistoryRepositoryAddRecordRequest struct {
	UserID    int       `json:"user_id"`
	A         geo.Point `json:"a"`
	B         geo.Point `json:"b"`
	Timestamp time.Time `json:"timestamp"`
}

// HistoryRepositoryGetDistanceRequest represents request object of HistoryRepository GetDistance method.
type HistoryRepositoryGetDistanceRequest struct {
	UserID int       `json:"user_id"`
	From   time.Time `json:"from"`
	To     time.Time `json:"to"`
}

// HistoryRepository represents history repository.
type HistoryRepository interface {
	AddRecord(ctx context.Context, req HistoryRepositoryAddRecordRequest) (domain.Record, error)
	GetDistance(ctx context.Context, req HistoryRepositoryGetDistanceRequest) (float64, error)
}
