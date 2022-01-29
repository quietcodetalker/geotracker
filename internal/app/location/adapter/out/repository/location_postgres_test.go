package repository_test

import (
	"context"
	"github.com/stretchr/testify/require"
	"gitlab.com/spacewalker/locations/internal/app/location/adapter/out/repository"
	"gitlab.com/spacewalker/locations/internal/app/location/core/domain"
	"gitlab.com/spacewalker/locations/internal/app/location/core/port"
	"gitlab.com/spacewalker/locations/internal/pkg/errpack"
	"gitlab.com/spacewalker/locations/internal/pkg/geo"
	"testing"
	"time"
)

func (s *PostgresTestSuite) Test_PostgresQueries_SetLocation() {
	createUserArgs := []port.CreateUserArg{
		{
			Username: "user1",
		},
		{
			Username: "user2",
		},
	}
	users := s.seedUsers(createUserArgs)

	setLocationArgs := []port.LocationRepositorySetLocationRequest{
		{
			UserID: users[0].ID,
			Point:  geo.Point{1.0, 1.0},
		},
	}
	locations := s.seedLocations(setLocationArgs)

	testCases := []struct {
		name   string
		arg    port.LocationRepositorySetLocationRequest
		hasErr bool
		isErr  error
		asErr  error
		assert func(t *testing.T, user domain.Location, err error)
	}{
		{
			name: "OK_UserExist_LocationExist",
			arg: port.LocationRepositorySetLocationRequest{
				UserID: users[0].ID,
				Point: geo.Point{
					locations[0].Point[0] + 1.0,
					locations[0].Point[1] + 1.0,
				},
			},
			hasErr: false,
			isErr:  nil,
			asErr:  nil,
			assert: func(t *testing.T, location domain.Location, err error) {
				require.Equal(t, users[0].ID, location.UserID)
				require.Equal(t, locations[0].Point.Latitude()+1.0, location.Point.Latitude())
				require.Equal(t, locations[0].Point.Longitude()+1.0, location.Point.Longitude())
				require.WithinDuration(t, locations[0].CreatedAt, location.CreatedAt, time.Second)
				require.WithinDuration(t, locations[0].UpdatedAt, location.UpdatedAt, time.Second)
			},
		},
		{
			name: "OK_UserExist_LocationDoesNotExist",
			arg: port.LocationRepositorySetLocationRequest{
				UserID: users[1].ID,
				Point:  geo.Point{1.0, 1.0},
			},
			hasErr: false,
			isErr:  nil,
			asErr:  nil,
			assert: func(t *testing.T, location domain.Location, err error) {
				require.Equal(t, users[1].ID, location.UserID)
				require.Equal(t, 1.0, location.Point.Latitude())
				require.Equal(t, 1.0, location.Point.Longitude())
				require.WithinDuration(t, time.Now(), location.CreatedAt, time.Second)
				require.WithinDuration(t, time.Now(), location.UpdatedAt, time.Second)
			},
		},
		{
			name: "ErrForeignKey_UserDoesNotExist_LocationDoesNotExist",
			arg: port.LocationRepositorySetLocationRequest{
				UserID: 0,
				Point:  geo.Point{1.0, 1.0},
			},
			hasErr: true,
			isErr:  nil,
			asErr:  nil,
			assert: func(t *testing.T, location domain.Location, err error) {
				require.Empty(t, location)
			},
		},
		{
			name: "ErrCheck_LattitudeLessThenMin",
			arg: port.LocationRepositorySetLocationRequest{
				UserID: users[0].ID,
				Point:  geo.Point{1.0, -181.0},
			},
			hasErr: true,
			isErr:  errpack.ErrInvalidArgument,
			assert: func(t *testing.T, location domain.Location, err error) {
				require.Empty(t, location)
			},
		},
		{
			name: "ErrCheck_LattitudeGreaterThenMax",
			arg: port.LocationRepositorySetLocationRequest{
				UserID: users[0].ID,
				Point:  geo.Point{1.0, 181.0},
			},
			hasErr: true,
			isErr:  errpack.ErrInvalidArgument,
			assert: func(t *testing.T, location domain.Location, err error) {
				require.Empty(t, location)
			},
		},
		{
			name: "ErrCheck_LongitudeLessThenMin",
			arg: port.LocationRepositorySetLocationRequest{
				UserID: users[0].ID,
				Point:  geo.Point{-181.0, 1.0},
			},
			hasErr: true,
			isErr:  errpack.ErrInvalidArgument,
			assert: func(t *testing.T, location domain.Location, err error) {
				require.Empty(t, location)
			},
		},
		{
			name: "ErrCheck_LongitudeGreaterThenMax",
			arg: port.LocationRepositorySetLocationRequest{
				UserID: users[0].ID,
				Point:  geo.Point{181.0, 1.0},
			},
			hasErr: true,
			isErr:  errpack.ErrInvalidArgument,
			assert: func(t *testing.T, location domain.Location, err error) {
				require.Empty(t, location)
			},
		},
	}

	repo := repository.NewPostgresRepository(s.db)

	for _, tc := range testCases {
		tc := tc
		s.T().Run(tc.name, func(t *testing.T) {
			location, err := repo.SetLocation(context.Background(), tc.arg)
			if !tc.hasErr {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				if tc.isErr != nil {
					require.ErrorIs(t, err, tc.isErr)
				}
				if tc.asErr != nil {
					require.ErrorAs(t, err, &tc.asErr)
				}
			}
			tc.assert(t, location, err)
		})
	}
}

func (s *PostgresTestSuite) Test_PostgresQueries_GetLocation() {
	createUserArgs := []port.CreateUserArg{
		{Username: "user0"},
		{Username: "user1"},
	}
	users := s.seedUsers(createUserArgs)

	setLocationRequests := []port.LocationRepositorySetLocationRequest{
		{
			UserID: users[0].ID,
			Point:  geo.Point{1.1, 2.2},
		},
	}
	locations := s.seedLocations(setLocationRequests)

	testCases := []struct {
		name     string
		userID   int
		expected domain.Location
		hasError bool
		isError  error
		asError  error
	}{
		{
			name:     "OK",
			userID:   locations[0].UserID,
			expected: locations[0],
			hasError: false,
		},
		{
			name:     "NoLocation",
			userID:   users[1].ID,
			expected: domain.Location{},
			hasError: true,
			isError:  errpack.ErrNotFound,
		},
		{
			name:     "NoUser",
			userID:   2,
			expected: domain.Location{},
			hasError: true,
			isError:  errpack.ErrNotFound,
		},
	}

	repo := repository.NewPostgresRepository(s.db)

	for _, tc := range testCases {
		tc := tc
		s.T().Run(tc.name, func(t *testing.T) {
			location, err := repo.GetLocation(context.Background(), tc.userID)
			if tc.hasError {
				require.Empty(t, location)
				if tc.isError != nil {
					require.ErrorIs(t, err, tc.isError)
				}
				if tc.asError != nil {
					require.ErrorAs(t, err, &tc.asError)
				}
			} else {
				require.Equal(t, tc.expected, location)
				require.NoError(t, err)
			}
		})
	}
}
