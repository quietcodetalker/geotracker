package repository_test

import (
	"context"
	"github.com/stretchr/testify/require"
	"gitlab.com/spacewalker/locations/internal/app/history/adapter/out/repository"
	"gitlab.com/spacewalker/locations/internal/app/history/core/domain"
	"gitlab.com/spacewalker/locations/internal/app/history/core/port"
	"gitlab.com/spacewalker/locations/internal/pkg/errpack"
	"gitlab.com/spacewalker/locations/internal/pkg/geo"
	"testing"
	"time"
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
			name: "OK_PointsTruncatedToFixedPrecision",
			req: port.HistoryRepositoryAddRecordRequest{
				UserID: 4,
				A:      geo.Point{0.123456789, 0.123456789},
				B:      geo.Point{0.123456789, 0.123456789},
			},
			assert: func(t *testing.T, rec domain.Record, err error) {
				require.NoError(t, err)
				require.NotEmpty(t, rec.ID)

				require.Equal(t, 0.12345678, rec.A.Longitude())
				require.Equal(t, 0.12345678, rec.A.Latitude())
				require.Equal(t, 0.12345678, rec.B.Longitude())
				require.Equal(t, 0.12345678, rec.B.Latitude())
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
				require.ErrorIs(t, err, errpack.ErrInvalidArgument)
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
				require.ErrorIs(t, err, errpack.ErrInvalidArgument)
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
				require.ErrorIs(t, err, errpack.ErrInvalidArgument)
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
				require.ErrorIs(t, err, errpack.ErrInvalidArgument)
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
				require.ErrorIs(t, err, errpack.ErrInvalidArgument)
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
				require.ErrorIs(t, err, errpack.ErrInvalidArgument)
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
				require.ErrorIs(t, err, errpack.ErrInvalidArgument)
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
				require.ErrorIs(t, err, errpack.ErrInvalidArgument)
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

func (s *PostgresTestSuite) Test_PostgresRepository_GetDistance() {
	ref := time.Now()
	records := []domain.Record{
		{
			UserID:    1,
			A:         geo.Point{0.0, 0.0},
			B:         geo.Point{1.0, 0.0},
			Timestamp: ref.Add(-time.Hour * 5),
		},
		{
			UserID:    1,
			A:         geo.Point{1.0, 0.0},
			B:         geo.Point{1.0, 1.0},
			Timestamp: ref.Add(-time.Hour * 4),
		},
		{
			UserID:    2,
			A:         geo.Point{1.0, 1.0},
			B:         geo.Point{2.0, 1.0},
			Timestamp: ref.Add(-time.Hour * 3),
		},
		{
			UserID:    1,
			A:         geo.Point{2.0, 1.0},
			B:         geo.Point{2.0, 2.0},
			Timestamp: ref.Add(-time.Hour * 2),
		},
	}

	s.seedRecords(records)

	testCases := []struct {
		name   string
		req    port.HistoryRepositoryGetDistanceRequest
		assert func(t *testing.T, distance float64, err error)
	}{
		{
			name: "OK_NoRecords",
			req: port.HistoryRepositoryGetDistanceRequest{
				UserID: 3,
				From:   ref.Add(-10 * time.Hour),
				To:     ref.Add(10 * time.Hour),
			},
			assert: func(t *testing.T, distance float64, err error) {
				require.NoError(t, err)
				require.Equal(t, 0.0, distance)
			},
		},
		{
			name: "OK",
			req: port.HistoryRepositoryGetDistanceRequest{
				UserID: 1,
				From:   ref.Add(-10 * time.Hour),
				To:     ref.Add(10 * time.Hour),
			},
			assert: func(t *testing.T, distance float64, err error) {
				require.NoError(t, err)
				require.InDelta(t, 111194.6977316823*3, distance, 0.00000001)
			},
		},
		{
			name: "OK_TwoRecordsInTimeFrame",
			req: port.HistoryRepositoryGetDistanceRequest{
				UserID: 1,
				From:   ref.Add(-10 * time.Hour),
				To:     ref.Add(-3 * time.Hour),
			},
			assert: func(t *testing.T, distance float64, err error) {
				require.NoError(t, err)
				require.InDelta(t, 111194.6977316823*2, distance, 0.00000001)
			},
		},
		{
			name: "OK_NoRecordsInTimeFrame",
			req: port.HistoryRepositoryGetDistanceRequest{
				UserID: 1,
				From:   ref.Add(-10 * time.Hour),
				To:     ref.Add(-6 * time.Hour),
			},
			assert: func(t *testing.T, distance float64, err error) {
				require.NoError(t, err)
				require.Equal(t, 0.0, distance)
			},
		},
	}

	repo := repository.NewPostgresRepository(s.db)

	for _, tc := range testCases {
		tc := tc
		s.T().Run(tc.name, func(t *testing.T) {
			distance, err := repo.GetDistance(context.Background(), tc.req)
			tc.assert(s.T(), distance, err)
		})
	}
}
