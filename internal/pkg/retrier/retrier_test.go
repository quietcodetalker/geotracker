//go:build utils
// +build utils

package retrier_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"gitlab.com/spacewalker/geotracker/internal/pkg/retrier"
)

type recorder struct {
	times   int
	counter int
	result  interface{}
	err     error
}

func (r *recorder) do(fn func()) func() (interface{}, error) {
	return func() (interface{}, error) {
		fn()
		defer func() {
			r.counter++
		}()
		if r.counter < r.times-1 {
			return nil, r.err
		}
		return r.result, nil
	}
}

func TestRetrier_Exec(t *testing.T) {
	testCases := []struct {
		name               string
		config             retrier.Config
		recorder           *recorder
		ctxTimeout         time.Duration
		expectedTimestamps []time.Time
		expectedRes        interface{}
		expectedErr        error
	}{
		{
			name: "0",
			config: retrier.Config{
				Delay:        3 * time.Second,
				Retries:      3,
				IsSuccessful: nil,
			},
			recorder: &recorder{
				times:  3,
				result: "test0",
				err:    errors.New("test0"),
			},
			ctxTimeout: time.Hour,
			expectedTimestamps: []time.Time{
				time.Now(),
				time.Now().Add(3 * time.Second),
				time.Now().Add(6 * time.Second),
			},
			expectedRes: "test0",
			expectedErr: nil,
		},
		{
			name: "1",
			config: retrier.Config{
				Delay:        3 * time.Second,
				Retries:      3,
				IsSuccessful: nil,
			},
			recorder: &recorder{
				times:  4,
				result: "test1",
				err:    errors.New("test1"),
			},
			ctxTimeout: time.Hour,
			expectedTimestamps: []time.Time{
				time.Now(),
				time.Now().Add(3 * time.Second),
				time.Now().Add(6 * time.Second),
			},
			expectedRes: nil,
			expectedErr: errors.New("test1"),
		},
		{
			name: "2",
			config: retrier.Config{
				Delay:        3 * time.Second,
				Retries:      3,
				IsSuccessful: nil,
			},
			recorder: &recorder{
				times:  2,
				result: "test2",
				err:    errors.New("test2"),
			},
			ctxTimeout: time.Hour,
			expectedTimestamps: []time.Time{
				time.Now(),
				time.Now().Add(3 * time.Second),
			},
			expectedRes: "test2",
			expectedErr: nil,
		},
		{
			name: "3 custom isSuccessful",
			config: retrier.Config{
				Delay:   3 * time.Second,
				Retries: 3,
				IsSuccessful: func(err error) bool {
					return err == nil || err.Error() == "invalid_argument"
				},
			},
			recorder: &recorder{
				times:  3,
				result: "test3",
				err:    errors.New("invalid_argument"),
			},
			ctxTimeout: time.Hour,
			expectedTimestamps: []time.Time{
				time.Now(),
			},
			expectedRes: nil,
			expectedErr: errors.New("invalid_argument"),
		},
		{
			name: "4 conetxt deadline exceeded",
			config: retrier.Config{
				Delay:   3 * time.Second,
				Retries: 3,
			},
			recorder: &recorder{
				times:  3,
				result: "test3",
				err:    errors.New("invalid_argument"),
			},
			ctxTimeout: time.Second,
			expectedTimestamps: []time.Time{
				time.Now(),
			},
			expectedRes: nil,
			expectedErr: context.DeadlineExceeded,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			r := retrier.New(tc.config)

			expected := tc.expectedTimestamps
			results := []time.Time{}
			fn := func() {
				results = append(results, time.Now())
			}
			rec := tc.recorder

			var wg sync.WaitGroup
			wg.Add(1)
			var res interface{}
			var err error

			ctx, cancel := context.WithTimeout(context.Background(), tc.ctxTimeout)
			defer cancel()

			go func() {
				res, err = r.Exec(ctx, rec.do(fn))
				wg.Done()
			}()
			wg.Wait()

			require.Equal(t, tc.expectedRes, res)
			require.Equal(t, tc.expectedErr, err)

			require.Len(t, results, len(expected))
			for i := 0; i < len(results); i++ {
				require.WithinDuration(t, expected[i], results[i], time.Second)
			}
		})
	}
}
