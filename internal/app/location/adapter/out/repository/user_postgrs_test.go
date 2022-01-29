package repository_test

import (
	"context"
	"github.com/stretchr/testify/require"
	"gitlab.com/spacewalker/locations/internal/app/location/adapter/out/repository"
	"gitlab.com/spacewalker/locations/internal/app/location/core/domain"
	"gitlab.com/spacewalker/locations/internal/app/location/core/port"
	"gitlab.com/spacewalker/locations/internal/pkg/errpack"
	"gitlab.com/spacewalker/locations/internal/pkg/geo"
	"gitlab.com/spacewalker/locations/internal/pkg/util"
	"testing"
	"time"
)

func (s *PostgresTestSuite) Test_PostgresQueries_CreateUser() {
	createUserArgs := []port.CreateUserArg{
		{
			Username: "user1",
		},
		{
			Username: "user2",
		},
	}

	users := s.seedUsers(createUserArgs)

	testCases := []struct {
		name   string
		arg    port.CreateUserArg
		hasErr bool
		isErr  error
		asErr  error
		assert func(t *testing.T, user domain.User, err error)
	}{
		{
			name:   "OK",
			arg:    port.CreateUserArg{Username: "user3"},
			hasErr: false,
			assert: func(t *testing.T, user domain.User, err error) {
				require.Equal(t, "user3", user.Username)
				for _, u := range users {
					require.NotEqual(t, u.ID, user.ID)
				}
			},
		},
		{
			name:   "ErrConstraint_UserAlreadyExists",
			arg:    port.CreateUserArg{Username: users[0].Username},
			hasErr: true,
			isErr:  errpack.ErrAlreadyExists,
			assert: func(t *testing.T, user domain.User, err error) {
				require.Empty(t, user)
			},
		},
		{
			name:   "ErrConstraint_InvalidUsername_TooShort",
			arg:    port.CreateUserArg{Username: "u"},
			hasErr: true,
			isErr:  errpack.ErrInvalidArgument,
			assert: func(t *testing.T, user domain.User, err error) {
				require.Empty(t, user)
			},
		},
		{
			name:   "ErrConstraint_InvalidUsername_TooLong",
			arg:    port.CreateUserArg{Username: "uuuuuuuuuuuuuuuuu"},
			hasErr: true,
			isErr:  errpack.ErrInvalidArgument,
			assert: func(t *testing.T, user domain.User, err error) {
				require.Empty(t, user)
			},
		},
		{
			name:   "ErrConstraint_InvalidUsername_DoesNotMatchPattern",
			arg:    port.CreateUserArg{Username: "user3_"},
			hasErr: true,
			isErr:  errpack.ErrInvalidArgument,
			assert: func(t *testing.T, user domain.User, err error) {
				require.Empty(t, user)
			},
		},
	}

	repo := repository.NewPostgresRepository(s.db)

	for _, tc := range testCases {
		tc := tc

		s.T().Run(tc.name, func(t *testing.T) {
			user, err := repo.CreateUser(context.Background(), tc.arg)

			if !tc.hasErr {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				if tc.isErr != nil {
					require.ErrorIs(t, err, tc.isErr)
				}
				if tc.asErr != nil {
					require.ErrorAs(t, err, tc.asErr)
				}
			}
			if tc.assert != nil {
				tc.assert(t, user, err)
			}
		})
	}
}

func (s *PostgresTestSuite) Test_PostgresQueries_GetUserByUsername() {
	createUserArgs := []port.CreateUserArg{
		{
			Username: "user1",
		},
		{
			Username: "user2",
		},
	}

	users := s.seedUsers(createUserArgs)

	testCases := []struct {
		name   string
		user   domain.User
		hasErr bool
		isErr  error
		asErr  error
		assert func(t *testing.T, user domain.User)
	}{
		{
			name:   "OK",
			user:   users[0],
			hasErr: false,
			isErr:  nil,
			asErr:  nil,
			assert: func(t *testing.T, user domain.User) {
				require.Equal(t, users[0].Username, user.Username)
				require.Equal(t, users[0].Username, user.Username)
				require.WithinDuration(t, users[0].CreatedAt, user.CreatedAt, time.Second)
				require.WithinDuration(t, users[0].UpdatedAt, user.UpdatedAt, time.Second)
			},
		},
		{
			name:   "NotFound",
			user:   domain.User{Username: "user3"},
			hasErr: true,
			isErr:  errpack.ErrNotFound,
			asErr:  nil,
			assert: func(t *testing.T, user domain.User) {
				require.Empty(t, user)
			},
		},
	}

	repo := repository.NewPostgresRepository(s.db)

	for _, tc := range testCases {
		tc := tc
		s.T().Run(tc.name, func(t *testing.T) {
			user, err := repo.GetByUsername(context.Background(), tc.user.Username)
			if !tc.hasErr {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				if tc.isErr != nil {
					require.ErrorIs(t, err, tc.isErr)
				}
				if tc.asErr != nil {
					require.ErrorAs(t, err, tc.asErr)
				}
			}
			if tc.assert != nil {
				tc.assert(t, user)
			}
		})
	}
}

func (s *PostgresTestSuite) Test_PostgresRepository_SetUserLocation() {
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
		arg    port.UserRepositorySetUserLocationRequest
		hasErr bool
		isErr  error
		asErr  error
		assert func(t *testing.T, res port.UserRepositorySetUserLocationResponse)
	}{
		{
			name: "OK_UserExist_LocationExist",
			arg: port.UserRepositorySetUserLocationRequest{
				Username: users[0].Username,
				Point:    geo.Point{locations[0].Point.Longitude() + 1.0, locations[0].Point.Latitude() + 1.0},
			},
			hasErr: false,
			isErr:  nil,
			asErr:  nil,
			assert: func(t *testing.T, res port.UserRepositorySetUserLocationResponse) {
				require.Equal(t, users[0].ID, res.Location.UserID)
				require.Equal(t, locations[0].Point.Latitude()+1.0, res.Location.Point.Latitude())
				require.Equal(t, locations[0].Point.Longitude()+1.0, res.Location.Point.Longitude())
				require.WithinDuration(t, locations[0].CreatedAt, res.Location.CreatedAt, time.Second)
				require.WithinDuration(t, locations[0].UpdatedAt, res.Location.UpdatedAt, time.Second)
			},
		},
		{
			name: "OK_UserExist_LocationDoesNotExist",
			arg: port.UserRepositorySetUserLocationRequest{
				Username: users[1].Username,
				Point:    geo.Point{1.0, 1.0},
			},
			hasErr: false,
			isErr:  nil,
			asErr:  nil,
			assert: func(t *testing.T, res port.UserRepositorySetUserLocationResponse) {
				require.Equal(t, users[1].ID, res.Location.UserID)
				require.Equal(t, 1.0, res.Location.Point.Latitude())
				require.Equal(t, 1.0, res.Location.Point.Longitude())
				require.WithinDuration(t, time.Now(), res.Location.CreatedAt, time.Second)
				require.WithinDuration(t, time.Now(), res.Location.UpdatedAt, time.Second)
			},
		},
		{
			name: "OK_UserDoesNotExist_LocationDoesNotExist",
			arg: port.UserRepositorySetUserLocationRequest{
				Username: "user3",
				Point:    geo.Point{1.0, 1.0},
			},
			hasErr: false,
			isErr:  nil,
			asErr:  nil,
			assert: func(t *testing.T, res port.UserRepositorySetUserLocationResponse) {
				for _, u := range users {
					require.NotEqual(t, u.ID, res.Location.UserID)
				}
				require.Equal(t, 1.0, res.Location.Point.Latitude())
				require.Equal(t, 1.0, res.Location.Point.Longitude())
				require.WithinDuration(t, time.Now(), res.Location.CreatedAt, time.Second)
				require.WithinDuration(t, time.Now(), res.Location.UpdatedAt, time.Second)
			},
		},
		{
			name: "OK_PointTruncatedToFixedPrecision",
			arg: port.UserRepositorySetUserLocationRequest{
				Username: "user3",
				Point:    geo.Point{0.123456789, 0.123456789},
			},
			hasErr: false,
			isErr:  nil,
			asErr:  nil,
			assert: func(t *testing.T, res port.UserRepositorySetUserLocationResponse) {
				require.Equal(t, 0.12345678, res.Location.Point.Latitude())
				require.Equal(t, 0.12345678, res.Location.Point.Longitude())
			},
		},
	}

	repo := repository.NewPostgresRepository(s.db)

	for _, tc := range testCases {
		tc := tc
		s.T().Run(tc.name, func(t *testing.T) {
			res, err := repo.SetUserLocation(context.Background(), tc.arg)
			if !tc.hasErr {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				if tc.isErr != nil {
					require.ErrorIs(t, err, tc.isErr)
				}
				if tc.asErr != nil {
					require.ErrorAs(t, err, tc.asErr)
				}
			}
			tc.assert(t, res)
		})
	}
}

func (s PostgresTestSuite) Test_PostgresQueries_ListUsersInRadius() {
	createUserArgs := []port.CreateUserArg{
		{Username: "user0"},
		{Username: "user1"},
		{Username: "user2"},
		{Username: "user3"},
		{Username: "user4"},
		{Username: "user5"},
		{Username: "user6"},
		{Username: "user7"},
		{Username: "user8"},
		{Username: "user9"},
		{Username: "user10"},
	}

	users := s.seedUsers(createUserArgs)

	setLocationArgs := []port.LocationRepositorySetLocationRequest{
		{
			UserID: users[0].ID,
			Point: geo.Point{
				0.0 + util.RandomFloat64(-0.01, 0.01),
				0.0 + util.RandomFloat64(-0.01, 0.01),
			},
		},
		{
			UserID: users[1].ID,
			Point: geo.Point{
				0.0 + util.RandomFloat64(-0.01, 0.01),
				-90.0 + util.RandomFloat64(0.0, 0.01),
			},
		},
		{
			UserID: users[2].ID,
			Point: geo.Point{
				0.0 + util.RandomFloat64(-0.01, 0.01),
				90.0 + util.RandomFloat64(-0.01, 0.0),
			},
		},
		{
			UserID: users[3].ID,
			Point: geo.Point{
				180.0 + util.RandomFloat64(-0.01, 0.0),
				0.0 + util.RandomFloat64(-0.01, 0.01),
			},
		},
		{
			UserID: users[4].ID,
			Point: geo.Point{
				-180.0 + util.RandomFloat64(0.0, 0.01),
				0.0 + util.RandomFloat64(-0.01, 0.01),
			},
		},
		{
			UserID: users[5].ID,
			Point: geo.Point{
				90.0 + util.RandomFloat64(-0.01, 0.01),
				0.0 + util.RandomFloat64(-0.01, 0.01),
			},
		},
		{
			UserID: users[6].ID,
			Point: geo.Point{
				-90.0 + util.RandomFloat64(-0.01, 0.01),
				0.0 + util.RandomFloat64(-0.01, 0.01),
			},
		},
		{
			UserID: users[7].ID,
			Point: geo.Point{
				180.0 + util.RandomFloat64(-0.01, 0.0),
				-90.0 + util.RandomFloat64(0.0, 0.01),
			},
		},
		{
			UserID: users[8].ID,
			Point: geo.Point{
				-180.0 + util.RandomFloat64(0.0, 0.01),
				-90.0 + util.RandomFloat64(0.0, 0.01),
			},
		},
		{
			UserID: users[9].ID,
			Point: geo.Point{
				180.0 + util.RandomFloat64(-0.01, 0.0),
				90.0 + util.RandomFloat64(-0.01, 0.0),
			},
		},
		{
			UserID: users[10].ID,
			Point: geo.Point{
				180.0 + util.RandomFloat64(-0.01, 0.0),
				-90.0 + util.RandomFloat64(0.0, 0.01),
			},
		},
	}

	s.seedLocations(setLocationArgs)

	testCases := []struct {
		name                  string
		in                    port.UserRepositoryListUsersInRadiusRequest
		hasErr                bool
		isErr                 error
		asErr                 error
		expectedUsers         []domain.User
		expectedNextPageToken int
	}{
		{
			name: "OK_0",
			in: port.UserRepositoryListUsersInRadiusRequest{
				Point:     geo.Point{0.0, 0.0},
				Radius:    10000.0,
				PageSize:  5,
				PageToken: 0,
			},
			hasErr:        false,
			expectedUsers: users[0:1],
		},
		{
			name: "OK_1",
			in: port.UserRepositoryListUsersInRadiusRequest{
				Point:     geo.Point{0.0, -90.0},
				Radius:    10000.0,
				PageSize:  5,
				PageToken: 0,
			},
			hasErr: false,
			expectedUsers: []domain.User{
				users[1],
				users[7],
				users[8],
				users[10],
			},
			expectedNextPageToken: 0,
		},
		{
			name: "OK_2",
			in: port.UserRepositoryListUsersInRadiusRequest{
				Point:     geo.Point{0.0, 90.0},
				Radius:    10000.0,
				PageSize:  5,
				PageToken: 0,
			},
			hasErr: false,
			expectedUsers: []domain.User{
				users[2],
				users[9],
			},
		},
		{
			name: "OK_3",
			in: port.UserRepositoryListUsersInRadiusRequest{
				Point:     geo.Point{180.0, 0.0},
				Radius:    10000.0,
				PageSize:  5,
				PageToken: 0,
			},
			hasErr: false,
			expectedUsers: []domain.User{
				users[3],
				users[4],
			},
		},
		{
			name: "OK_4",
			in: port.UserRepositoryListUsersInRadiusRequest{
				Point:     geo.Point{90.0, 0.0},
				Radius:    10000.0,
				PageSize:  5,
				PageToken: 0,
			},
			hasErr: false,
			expectedUsers: []domain.User{
				users[5],
			},
		},
		{
			name: "OK_5",
			in: port.UserRepositoryListUsersInRadiusRequest{
				Point:     geo.Point{-90.0, 0.0},
				Radius:    10000.0,
				PageSize:  5,
				PageToken: 0,
			},
			hasErr: false,
			expectedUsers: []domain.User{
				users[6],
			},
		},
		{
			name: "OK_6",
			in: port.UserRepositoryListUsersInRadiusRequest{
				Point:     geo.Point{45.0, 45.0},
				Radius:    10.0,
				PageSize:  5,
				PageToken: 0,
			},
			hasErr:                false,
			expectedUsers:         nil,
			expectedNextPageToken: 0,
		},
		{
			name: "OK_Pagination_1",
			in: port.UserRepositoryListUsersInRadiusRequest{
				Point:     geo.Point{0.0, -90.0},
				Radius:    10000.0,
				PageSize:  2,
				PageToken: 0,
			},
			hasErr: false,
			expectedUsers: []domain.User{
				users[1],
				users[7],
			},
			expectedNextPageToken: users[7].ID,
		},
		{
			name: "OK_Pagination_2",
			in: port.UserRepositoryListUsersInRadiusRequest{
				Point:     geo.Point{0.0, -90.0},
				Radius:    10000.0,
				PageSize:  1,
				PageToken: users[1].ID,
			},
			hasErr: false,
			expectedUsers: []domain.User{
				users[7],
			},
			expectedNextPageToken: users[7].ID,
		},
	}

	repo := repository.NewPostgresRepository(s.db)

	for _, tc := range testCases {
		tc := tc
		s.T().Run(tc.name, func(t *testing.T) {
			res, err := repo.ListUsersInRadius(context.Background(), tc.in)

			require.Equal(t, tc.expectedUsers, res.Users)
			require.Equal(t, tc.expectedNextPageToken, res.NextPageToken)

			if !tc.hasErr {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				if tc.isErr != nil {
					require.ErrorIs(t, err, tc.isErr)
				}
				if tc.asErr != nil {
					require.ErrorAs(t, err, tc.asErr)
				}
			}
		})
	}
}
