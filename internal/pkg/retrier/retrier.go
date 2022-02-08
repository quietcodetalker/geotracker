package retrier

import (
	"context"
	"time"
)

const (
	defaultDelay   = time.Second
	defaultRetries = 3
)

func defaultIsSuccessful(err error) bool {
	return err == nil
}

// Retrier implements retry with backoff patterns.
type Retrier struct {
	delay        time.Duration
	retries      int
	isSuccessful func(err error) bool
}

// Config is a retrier configuration structure.
type Config struct {
	Delay        time.Duration
	Retries      int
	IsSuccessful func(err error) bool
}

// New returns a pointer to new instance of Retrier.
func New(cfg Config) *Retrier {
	r := Retrier{
		delay:        cfg.Delay,
		retries:      cfg.Retries,
		isSuccessful: cfg.IsSuccessful,
	}
	if r.isSuccessful == nil {
		r.isSuccessful = defaultIsSuccessful
	}
	if r.delay == 0 {
		r.delay = defaultDelay
	}
	if r.retries <= 0 {
		r.retries = defaultRetries
	}

	return &r
}

// Exec executes provided function with retry and backoff patterns.
func (r *Retrier) Exec(ctx context.Context, fn func() (interface{}, error)) (interface{}, error) {
	for i := 0; ; i++ {
		res, err := fn()
		if r.isSuccessful(err) || i >= r.retries-1 {
			return res, err
		}

		select {
		case <-time.After(r.delay):
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}
