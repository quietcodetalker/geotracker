package historyclient

import (
	"context"
	"errors"
	"fmt"

	"gitlab.com/spacewalker/geotracker/internal/pkg/retrier"

	"gitlab.com/spacewalker/geotracker/internal/pkg/errpack"

	"github.com/sony/gobreaker"

	"gitlab.com/spacewalker/geotracker/internal/app/location/core/port"
)

// Proxy wraps history client and applies circuit breaker and retry with backoff patterns.
type Proxy struct {
	client  port.HistoryClient
	breaker *gobreaker.CircuitBreaker
	retrier *retrier.Retrier
}

// NewProxy returns a new instance of Proxy.
func NewProxy(
	client port.HistoryClient,
	breaker *gobreaker.CircuitBreaker,
	retrier *retrier.Retrier,
) port.HistoryClient {
	return &Proxy{
		client:  client,
		breaker: breaker,
		retrier: retrier,
	}
}

// AddRecord TODO: description
func (p *Proxy) AddRecord(ctx context.Context, req port.HistoryClientAddRecordRequest) (port.HistoryClientAddRecordResponse, error) {
	res, err := p.retrier.Exec(ctx, func() (interface{}, error) {
		res, err := p.breaker.Execute(func() (interface{}, error) {
			return p.client.AddRecord(ctx, req)
		})

		if errors.Is(err, gobreaker.ErrOpenState) || errors.Is(err, gobreaker.ErrTooManyRequests) {
			return res.(port.HistoryClientAddRecordResponse), fmt.Errorf("%w: %v", errpack.ErrInternalError, err)
		}

		return res, err
	})

	if err != nil {
		return port.HistoryClientAddRecordResponse{}, err
	}

	return res.(port.HistoryClientAddRecordResponse), nil
}
