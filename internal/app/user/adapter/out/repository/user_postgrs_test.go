package repository_test

import (
	"context"
	"github.com/stretchr/testify/require"
	"gitlab.com/spacewalker/locations/internal/app/user/adapter/out/repository"
	"gitlab.com/spacewalker/locations/internal/app/user/core/domain"
	"gitlab.com/spacewalker/locations/internal/app/user/core/port"
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
			isErr:  port.ErrAlreadyExists,
			assert: func(t *testing.T, user domain.User, err error) {
				require.Empty(t, user)
			},
		},
		{
			name:   "ErrConstraint_InvalidUsername_TooShort",
			arg:    port.CreateUserArg{Username: "u"},
			hasErr: true,
			isErr:  port.ErrInvalidUsername,
			assert: func(t *testing.T, user domain.User, err error) {
				require.Empty(t, user)
			},
		},
		{
			name:   "ErrConstraint_InvalidUsername_TooLong",
			arg:    port.CreateUserArg{Username: "u"},
			hasErr: true,
			isErr:  port.ErrInvalidUsername,
			assert: func(t *testing.T, user domain.User, err error) {
				require.Empty(t, user)
			},
		},
		{
			name:   "ErrConstraint_InvalidUsername_DoesNotMatchPattern",
			arg:    port.CreateUserArg{Username: "user3_"},
			hasErr: true,
			isErr:  port.ErrInvalidUsername,
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
			isErr:  port.ErrNotFound,
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

	setLocationArgs := []port.SetLocationArg{
		{
			UserID:    users[0].ID,
			Latitude:  1.0,
			Longitude: 1.0,
		},
	}
	locations := s.seedLocations(setLocationArgs)

	testCases := []struct {
		name   string
		arg    port.SetUserLocationArg
		hasErr bool
		isErr  error
		asErr  error
		assert func(t *testing.T, user domain.Location)
	}{
		{
			name: "OK_UserExist_LocationExist",
			arg: port.SetUserLocationArg{
				Username:  users[0].Username,
				Latitude:  locations[0].Latitude + 1.0,
				Longitude: locations[0].Longitude + 1.0,
			},
			hasErr: false,
			isErr:  nil,
			asErr:  nil,
			assert: func(t *testing.T, location domain.Location) {
				require.Equal(t, users[0].ID, location.UserID)
				require.Equal(t, locations[0].Latitude+1.0, location.Latitude)
				require.Equal(t, locations[0].Longitude+1.0, location.Longitude)
				require.WithinDuration(t, locations[0].CreatedAt, location.CreatedAt, time.Second)
				require.WithinDuration(t, locations[0].UpdatedAt, location.UpdatedAt, time.Second)
			},
		},
		{
			name: "OK_UserExist_LocationDoesNotExist",
			arg: port.SetUserLocationArg{
				Username:  users[1].Username,
				Latitude:  1.0,
				Longitude: 1.0,
			},
			hasErr: false,
			isErr:  nil,
			asErr:  nil,
			assert: func(t *testing.T, location domain.Location) {
				require.Equal(t, users[1].ID, location.UserID)
				require.Equal(t, 1.0, location.Latitude)
				require.Equal(t, 1.0, location.Longitude)
				require.WithinDuration(t, time.Now(), location.CreatedAt, time.Second)
				require.WithinDuration(t, time.Now(), location.UpdatedAt, time.Second)
			},
		},
		{
			name: "OK_UserDoesNotExist_LocationDoesNotExist",
			arg: port.SetUserLocationArg{
				Username:  "user3",
				Latitude:  1.0,
				Longitude: 1.0,
			},
			hasErr: false,
			isErr:  nil,
			asErr:  nil,
			assert: func(t *testing.T, location domain.Location) {
				for _, u := range users {
					require.NotEqual(t, u.ID, location.UserID)
				}
				require.Equal(t, 1.0, location.Latitude)
				require.Equal(t, 1.0, location.Longitude)
				require.WithinDuration(t, time.Now(), location.CreatedAt, time.Second)
				require.WithinDuration(t, time.Now(), location.UpdatedAt, time.Second)
			},
		},
	}

	repo := repository.NewPostgresRepository(s.db)

	for _, tc := range testCases {
		tc := tc
		s.T().Run(tc.name, func(t *testing.T) {
			location, err := repo.SetUserLocation(context.Background(), tc.arg)
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
			tc.assert(t, location)
		})
	}
}
