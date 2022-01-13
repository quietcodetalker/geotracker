package pagination_test

import (
	"github.com/stretchr/testify/require"
	"gitlab.com/spacewalker/locations/internal/pkg/util/pagination"
	"testing"
)

func TestEncodeCursor(t *testing.T) {
	testCases := []struct {
		name      string
		pageSize  int
		pageToken int
		expected  string
	}{
		{
			name:      "0 0",
			pageToken: 0,
			pageSize:  0,
			expected:  "MCAw",
		},
		{
			name:      "100 100",
			pageToken: 100,
			pageSize:  100,
			expected:  "MTAwIDEwMA==",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			cursor := pagination.EncodeCursor(tc.pageToken, tc.pageSize)
			require.Equal(t, tc.expected, cursor)
		})
	}
}

func TestDecodeCursor(t *testing.T) {
	testCases := []struct {
		name              string
		cursor            string
		expectedPageToken int
		expectedPageSize  int
		expectedErr       error
	}{
		{
			name:              "OK_0_0",
			cursor:            "MCAw",
			expectedPageToken: 0,
			expectedPageSize:  0,
			expectedErr:       nil,
		},
		{
			name:              "OK_100_100",
			cursor:            "MTAwIDEwMA==",
			expectedPageToken: 100,
			expectedPageSize:  100,
			expectedErr:       nil,
		},
		{
			name:              "InvalidFormat_10000",
			cursor:            "MTAwMTAw",
			expectedPageToken: 0,
			expectedPageSize:  0,
			expectedErr:       pagination.ErrInvalidCursor,
		},
		{
			name:              "InvalidBase64",
			cursor:            "invalid",
			expectedPageToken: 0,
			expectedPageSize:  0,
			expectedErr:       pagination.ErrInvalidCursor,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			pageToken, pageSize, err := pagination.DecodeCursor(tc.cursor)
			require.Equal(t, tc.expectedPageToken, pageToken)
			require.Equal(t, tc.expectedPageSize, pageSize)
			if tc.expectedErr == nil {
				require.NoError(t, err)
			} else {
				require.ErrorIs(t, err, tc.expectedErr)
			}
		})
	}
}
