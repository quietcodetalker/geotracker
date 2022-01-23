package util_test

import (
	"github.com/stretchr/testify/require"
	"gitlab.com/spacewalker/locations/internal/pkg/util"
	"testing"
)

func TestTrunc(t *testing.T) {
	testCases := []struct {
		name      string
		number    float64
		precision int
		expected  float64
	}{
		{
			name:      "OK_1",
			number:    0.123456789,
			precision: 8,
			expected:  0.12345678,
		},
		{
			name:      "OK_2",
			number:    0.1234,
			precision: 8,
			expected:  0.1234,
		},
		{
			name:      "OK_3",
			number:    150.1234,
			precision: -2,
			expected:  100.0,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := util.Trunc(tc.number, tc.precision)
			require.Equal(t, tc.expected, got)
		})
	}
}
