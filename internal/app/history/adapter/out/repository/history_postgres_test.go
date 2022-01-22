package repository_test

import (
	"context"
	"github.com/stretchr/testify/require"
	"gitlab.com/spacewalker/locations/internal/app/history/adapter/out/repository"
	"gitlab.com/spacewalker/locations/internal/app/history/core/domain"
	"gitlab.com/spacewalker/locations/internal/app/history/core/port"
	"gitlab.com/spacewalker/locations/internal/pkg/geo"
	"testing"
)

func (s *PostgresTestSuite) Test_PostgresRepository_AddRecord() {
	testCases := []struct {
		name   string
		req    port.HistoryRepositoryAddRecordRequest
		assert func(t *testing.T, rec domain.Record, err error)
	}{
		{
			name: "OK",
			req: port.HistoryRepositoryAddRecordRequest{
				UserID: 1,
				A:      geo.Point{0, 0},
				B:      geo.Point{1, 1},
			},
			assert: func(t *testing.T, rec domain.Record, err error) {
				require.NoError(t, err)
				require.NotEmpty(t, rec.ID)

				require.Equal(t, 1, rec.UserID)
				require.Equal(t, 0.0, rec.A.Longitude())
				require.Equal(t, 0.0, rec.A.Latitude())
				require.Equal(t, 1.0, rec.B.Longitude())
				require.Equal(t, 1.0, rec.B.Latitude())
			},
		},
		{
			name: "InvalidPoint_A_Longitude_LessThanMin",
			req: port.HistoryRepositoryAddRecordRequest{
				UserID: 1,
				A:      geo.Point{-180.1, 0},
				B:      geo.Point{1, 1},
			},
			assert: func(t *testing.T, rec domain.Record, err error) {
				require.Empty(t, rec)

				var invalidLocationError *port.InvalidLocationError
				require.ErrorAs(t, err, &invalidLocationError)
				require.Equal(t, port.InvalidLocationError{
					Violations: []port.InvalidLocationErrorViolation{
						{
							Subject: "a.longitude",
							Value:   -180.1,
						},
					},
				}, *invalidLocationError)
			},
		},
		{
			name: "InvalidPoint_A_Longitude_GreaterThanMax",
			req: port.HistoryRepositoryAddRecordRequest{
				UserID: 1,
				A:      geo.Point{180.1, 0},
				B:      geo.Point{1, 1},
			},
			assert: func(t *testing.T, rec domain.Record, err error) {
				require.Empty(t, rec)

				var invalidLocationError *port.InvalidLocationError
				require.ErrorAs(t, err, &invalidLocationError)
				require.Equal(t, port.InvalidLocationError{
					Violations: []port.InvalidLocationErrorViolation{
						{
							Subject: "a.longitude",
							Value:   180.1,
						},
					},
				}, *invalidLocationError)
			},
		},
		{
			name: "InvalidPoint_A_Latitude_LessThanMin",
			req: port.HistoryRepositoryAddRecordRequest{
				UserID: 1,
				A:      geo.Point{0, -90.1},
				B:      geo.Point{1, 1},
			},
			assert: func(t *testing.T, rec domain.Record, err error) {
				require.Empty(t, rec)

				var invalidLocationError *port.InvalidLocationError
				require.ErrorAs(t, err, &invalidLocationError)
				require.Equal(t, port.InvalidLocationError{
					Violations: []port.InvalidLocationErrorViolation{
						{
							Subject: "a.latitude",
							Value:   -90.1,
						},
					},
				}, *invalidLocationError)
			},
		},
		{
			name: "InvalidPoint_A_Latitude_GreaterThanMax",
			req: port.HistoryRepositoryAddRecordRequest{
				UserID: 1,
				A:      geo.Point{0, 90.1},
				B:      geo.Point{1, 1},
			},
			assert: func(t *testing.T, rec domain.Record, err error) {
				require.Empty(t, rec)

				var invalidLocationError *port.InvalidLocationError
				require.ErrorAs(t, err, &invalidLocationError)
				require.Equal(t, port.InvalidLocationError{
					Violations: []port.InvalidLocationErrorViolation{
						{
							Subject: "a.latitude",
							Value:   90.1,
						},
					},
				}, *invalidLocationError)
			},
		},
		{
			name: "InvalidPoint_B_Longitude_LessThanMin",
			req: port.HistoryRepositoryAddRecordRequest{
				UserID: 1,
				A:      geo.Point{1, 1},
				B:      geo.Point{-180.1, 0},
			},
			assert: func(t *testing.T, rec domain.Record, err error) {
				require.Empty(t, rec)

				var invalidLocationError *port.InvalidLocationError
				require.ErrorAs(t, err, &invalidLocationError)
				require.Equal(t, port.InvalidLocationError{
					Violations: []port.InvalidLocationErrorViolation{
						{
							Subject: "b.longitude",
							Value:   -180.1,
						},
					},
				}, *invalidLocationError)
			},
		},
		{
			name: "InvalidPoint_B_Longitude_GreaterThanMax",
			req: port.HistoryRepositoryAddRecordRequest{
				UserID: 1,
				A:      geo.Point{1, 1},
				B:      geo.Point{180.1, 0},
			},
			assert: func(t *testing.T, rec domain.Record, err error) {
				require.Empty(t, rec)

				var invalidLocationError *port.InvalidLocationError
				require.ErrorAs(t, err, &invalidLocationError)
				require.Equal(t, port.InvalidLocationError{
					Violations: []port.InvalidLocationErrorViolation{
						{
							Subject: "b.longitude",
							Value:   180.1,
						},
					},
				}, *invalidLocationError)
			},
		},
		{
			name: "InvalidPoint_B_Latitude_LessThanMin",
			req: port.HistoryRepositoryAddRecordRequest{
				UserID: 1,
				A:      geo.Point{1, 1},
				B:      geo.Point{0, -90.1},
			},
			assert: func(t *testing.T, rec domain.Record, err error) {
				require.Empty(t, rec)

				var invalidLocationError *port.InvalidLocationError
				require.ErrorAs(t, err, &invalidLocationError)
				require.Equal(t, port.InvalidLocationError{
					Violations: []port.InvalidLocationErrorViolation{
						{
							Subject: "b.latitude",
							Value:   -90.1,
						},
					},
				}, *invalidLocationError)
			},
		},
		{
			name: "InvalidPoint_B_Latitude_GreaterThanMax",
			req: port.HistoryRepositoryAddRecordRequest{
				UserID: 1,
				A:      geo.Point{0, 90.1},
				B:      geo.Point{1, 1},
			},
			assert: func(t *testing.T, rec domain.Record, err error) {
				require.Empty(t, rec)

				var invalidLocationError *port.InvalidLocationError
				require.ErrorAs(t, err, &invalidLocationError)
				require.Equal(t, port.InvalidLocationError{
					Violations: []port.InvalidLocationErrorViolation{
						{
							Subject: "a.latitude",
							Value:   90.1,
						},
					},
				}, *invalidLocationError)
			},
		},
	}

	repo := repository.NewPostgresRepository(s.db)

	for _, tc := range testCases {
		tc := tc
		s.T().Run(tc.name, func(t *testing.T) {
			rec, err := repo.AddRecord(context.Background(), tc.req)
			tc.assert(s.T(), rec, err)
		})
	}
}
