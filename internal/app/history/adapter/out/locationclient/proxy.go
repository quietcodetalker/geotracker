package locationclient

import (
	"context"
	"errors"
	"fmt"

	"gitlab.com/spacewalker/locations/internal/app/history/core/port"

	"gitlab.com/spacewalker/locations/internal/pkg/retrier"

	"gitlab.com/spacewalker/locations/internal/pkg/errpack"

	"github.com/sony/gobreaker"
)

// Proxy wraps history client and applies circuit breaker and retry with backoff patterns.
type Proxy struct {
	client  port.LocationClient
	breaker *gobreaker.CircuitBreaker
	retrier *retrier.Retrier
}

// NewProxy returns a new instance of Proxy.
func NewProxy(
	client port.LocationClient,
	breaker *gobreaker.CircuitBreaker,
	retrier *retrier.Retrier,
) port.LocationClient {
	return &Proxy{
		client:  client,
		breaker: breaker,
		retrier: retrier,
	}
}

func (p *Proxy) GetUserIDByUsername(ctx context.Context, username string) (int, error) {
	res, err := p.retrier.Exec(ctx, func() (interface{}, error) {
		res, err := p.breaker.Execute(func() (interface{}, error) {
			return p.client.GetUserIDByUsername(ctx, username)
		})

		if errors.Is(err, gobreaker.ErrOpenState) || errors.Is(err, gobreaker.ErrTooManyRequests) {
			return res, fmt.Errorf("%w: %v", errpack.ErrInternalError, err)
		}

		return res, err
	})

	if err != nil {
		return 0, err
	}

	return res.(int), nil
}
