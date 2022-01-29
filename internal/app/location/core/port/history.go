//go:generate mockgen -destination=mock/mock_history.go -package=mock . HistoryClient
package port

import (
	"context"
	"gitlab.com/spacewalker/locations/internal/pkg/geo"
	"time"
)

type HistoryClientAddRecordRequest struct {
	UserID    int       `json:"user_id"`
	A         geo.Point `json:"a"`
	B         geo.Point `json:"b"`
	Timestamp time.Time `json:"timestamp"`
}

type HistoryClientAddRecordResponse struct {
	UserID    int       `json:"user_id"`
	A         geo.Point `json:"a"`
	B         geo.Point `json:"b"`
	Timestamp time.Time `json:"timestamp"`
}

type HistoryClientGetDistanceRequest struct{}
type HistoryClientGetDistanceResponse struct{}

type HistoryClient interface {
	AddRecord(ctx context.Context, req HistoryClientAddRecordRequest) (HistoryClientAddRecordResponse, error)
	GetDistance(ctx context.Context, req HistoryClientGetDistanceRequest) (HistoryClientGetDistanceResponse, error)
}
