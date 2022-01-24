//go:generate mockgen -destination=mock/mock_history.go -package=mock . HistoryClient
package port

import (
	"context"
	"gitlab.com/spacewalker/locations/internal/pkg/geo"
	"time"
)

type HistoryClientAddRecordRequest struct {
	UserID int       `json:"user_id"`
	A      geo.Point `json:"a"`
	B      geo.Point `json:"b"`
}

type HistoryClientAddRecordResponse struct {
	UserID    int       `json:"user_id"`
	A         geo.Point `json:"a"`
	B         geo.Point `json:"b"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type HistoryClientGetDistanceRequest struct{}
type HistoryClientGetDistanceResponse struct{}

type HistoryClient interface {
	AddRecord(ctx context.Context, req HistoryClientAddRecordRequest) (HistoryClientAddRecordResponse, error)
	GetDistance(ctx context.Context, req HistoryClientGetDistanceRequest) (HistoryClientGetDistanceResponse, error)
}
