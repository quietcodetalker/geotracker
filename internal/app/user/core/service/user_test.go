package service_test

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gitlab.com/spacewalker/locations/internal/app/user/core/domain"
	"gitlab.com/spacewalker/locations/internal/app/user/core/port"
	"gitlab.com/spacewalker/locations/internal/app/user/core/port/mock"
	"gitlab.com/spacewalker/locations/internal/app/user/core/service"
	"gitlab.com/spacewalker/locations/internal/pkg/util"
	"testing"
	"time"
)

type UserSvcTestSuite struct {
	suite.Suite
}

func (s *UserSvcTestSuite) SetupTest() {}

func (s *UserSvcTestSuite) TearDownTest() {}

func (s *UserSvcTestSuite) TearDownSuite() {}

func TestSvcTestSuite(t *testing.T) {
	// Skip tests when using "-short" flag.
	if testing.Short() {
		t.Skip("Skipping long-running tests")
	}

	suite.Run(t, new(UserSvcTestSuite))
}

func (s *UserSvcTestSuite) Test_UserService_SetUserLocation() {
	username := "user1"

	point := domain.Point{
		util.RandomFloat64(-180.0, 180.0),
		util.RandomFloat64(-90.0, 90.0),
	}

	testCases := []struct {
		name      string
		arg       port.SetUserLocationRequest
		buildStus func(repo *mock.MockUserRepository)
		assert    func(t *testing.T, res domain.Location, err error)
	}{
		{
			name: "OK",
			arg: port.SetUserLocationRequest{
				Username:  username,
				Longitude: point.Longitude(),
				Latitude:  point.Latitude(),
			},
			buildStus: func(repo *mock.MockUserRepository) {
				repo.EXPECT().
					SetUserLocation(
						gomock.Any(),
						gomock.Eq(port.SetUserLocationArg{
							Username: username,
							Point:    point,
						}),
					).
					Times(1).
					Return(
						domain.Location{
							UserID:    1,
							Point:     point,
							CreatedAt: time.Now(),
							UpdatedAt: time.Now(),
						},
						nil,
					)
			},
			assert: func(t *testing.T, res domain.Location, err error) {
				require.NoError(t, err)
				require.Equal(t, 1, res.UserID)
				require.Equal(t, point.Longitude(), res.Point.Longitude())
				require.Equal(t, point.Latitude(), res.Point.Latitude())

				require.WithinDuration(t, time.Now(), res.CreatedAt, time.Second)
				require.WithinDuration(t, time.Now(), res.UpdatedAt, time.Second)
			},
		},
		{
			name: "OK",
			arg: port.SetUserLocationRequest{
				Username:  username,
				Longitude: -180.0,
				Latitude:  -90.0,
			},
			buildStus: func(repo *mock.MockUserRepository) {
				repo.EXPECT().
					SetUserLocation(
						gomock.Any(),
						gomock.Eq(port.SetUserLocationArg{
							Username: username,
							Point:    domain.Point{-180.0, -90.0},
						}),
					).
					Times(1).
					Return(
						domain.Location{
							UserID:    1,
							Point:     domain.Point{-180.0, -90.0},
							CreatedAt: time.Now(),
							UpdatedAt: time.Now(),
						},
						nil,
					)
			},
			assert: func(t *testing.T, res domain.Location, err error) {
				require.NoError(t, err)
				require.Equal(t, 1, res.UserID)
				require.Equal(t, point.Longitude(), res.Point.Longitude())
				require.Equal(t, point.Latitude(), res.Point.Latitude())

				require.WithinDuration(t, time.Now(), res.CreatedAt, time.Second)
				require.WithinDuration(t, time.Now(), res.UpdatedAt, time.Second)
			},
		},
		{
			name: "OK",
			arg: port.SetUserLocationRequest{
				Username:  username,
				Longitude: 180.0,
				Latitude:  90.0,
			},
			buildStus: func(repo *mock.MockUserRepository) {
				repo.EXPECT().
					SetUserLocation(
						gomock.Any(),
						gomock.Eq(port.SetUserLocationArg{
							Username: username,
							Point:    domain.Point{180.0, 90.0},
						}),
					).
					Times(1).
					Return(
						domain.Location{
							UserID:    1,
							Point:     domain.Point{180.0, 90.0},
							CreatedAt: time.Now(),
							UpdatedAt: time.Now(),
						},
						nil,
					)
			},
			assert: func(t *testing.T, res domain.Location, err error) {
				require.NoError(t, err)
				require.Equal(t, 1, res.UserID)
				require.Equal(t, point.Longitude(), res.Point.Longitude())
				require.Equal(t, point.Latitude(), res.Point.Latitude())

				require.WithinDuration(t, time.Now(), res.CreatedAt, time.Second)
				require.WithinDuration(t, time.Now(), res.UpdatedAt, time.Second)
			},
		},
		{
			name: "OK",
			arg: port.SetUserLocationRequest{
				Username:  username,
				Longitude: 0.0,
				Latitude:  0.0,
			},
			buildStus: func(repo *mock.MockUserRepository) {
				repo.EXPECT().
					SetUserLocation(
						gomock.Any(),
						gomock.Eq(port.SetUserLocationArg{
							Username: username,
							Point:    domain.Point{0, 0},
						}),
					).
					Times(1).
					Return(
						domain.Location{
							UserID:    1,
							Point:     domain.Point{0, 0},
							CreatedAt: time.Now(),
							UpdatedAt: time.Now(),
						},
						nil,
					)
			},
			assert: func(t *testing.T, res domain.Location, err error) {
				require.NoError(t, err)
				require.Equal(t, 1, res.UserID)
				require.Equal(t, point.Longitude(), res.Point.Longitude())
				require.Equal(t, point.Latitude(), res.Point.Latitude())

				require.WithinDuration(t, time.Now(), res.CreatedAt, time.Second)
				require.WithinDuration(t, time.Now(), res.UpdatedAt, time.Second)
			},
		},
		{
			name: "Err_InvalidLongitude_LessThanMin",
			arg: port.SetUserLocationRequest{
				Username:  username,
				Longitude: -180.01,
				Latitude:  point.Latitude(),
			},
			buildStus: func(repo *mock.MockUserRepository) {
				repo.EXPECT().
					SetUserLocation(
						gomock.Any(),
						gomock.Any(),
					).
					Times(0)
			},
			assert: func(t *testing.T, res domain.Location, err error) {
				require.Empty(t, res)
				var invalidArgumentErr *port.InvalidArgumentError
				require.ErrorAs(t, err, &invalidArgumentErr)
			},
		},
		{
			name: "Err_InvalidLongitude_GreaterThanMax",
			arg: port.SetUserLocationRequest{
				Username:  username,
				Longitude: 180.01,
				Latitude:  point.Latitude(),
			},
			buildStus: func(repo *mock.MockUserRepository) {
				repo.EXPECT().
					SetUserLocation(
						gomock.Any(),
						gomock.Any(),
					).
					Times(0)
			},
			assert: func(t *testing.T, res domain.Location, err error) {
				require.Empty(t, res)
				var invalidArgumentErr *port.InvalidArgumentError
				require.ErrorAs(t, err, &invalidArgumentErr)
			},
		},
		{
			name: "Err_InvalidLatitude_LessThanMin",
			arg: port.SetUserLocationRequest{
				Username:  username,
				Longitude: point.Longitude(),
				Latitude:  -90.01,
			},
			buildStus: func(repo *mock.MockUserRepository) {
				repo.EXPECT().
					SetUserLocation(
						gomock.Any(),
						gomock.Any(),
					).
					Times(0)
			},
			assert: func(t *testing.T, res domain.Location, err error) {
				require.Empty(t, res)
				var invalidArgumentErr *port.InvalidArgumentError
				require.ErrorAs(t, err, &invalidArgumentErr)
			},
		},
		{
			name: "Err_InvalidLatitude_GreaterThanMax",
			arg: port.SetUserLocationRequest{
				Username:  username,
				Longitude: point.Longitude(),
				Latitude:  90.01,
			},
			buildStus: func(repo *mock.MockUserRepository) {
				repo.EXPECT().
					SetUserLocation(
						gomock.Any(),
						gomock.Any(),
					).
					Times(0)
			},
			assert: func(t *testing.T, res domain.Location, err error) {
				require.Empty(t, res)
				var invalidArgumentErr *port.InvalidArgumentError
				require.ErrorAs(t, err, &invalidArgumentErr)
			},
		},
		{
			name: "Err_InvalidUsername_TooShort",
			arg: port.SetUserLocationRequest{
				Username:  "u",
				Longitude: point.Longitude(),
				Latitude:  point.Latitude(),
			},
			buildStus: func(repo *mock.MockUserRepository) {
				repo.EXPECT().
					SetUserLocation(
						gomock.Any(),
						gomock.Any(),
					).
					Times(0)
			},
			assert: func(t *testing.T, res domain.Location, err error) {
				require.Empty(t, res)
				var invalidArgumentErr *port.InvalidArgumentError
				require.ErrorAs(t, err, &invalidArgumentErr)
			},
		},
		{
			name: "Err_InvalidUsername_TooLong",
			arg: port.SetUserLocationRequest{
				Username:  "uuuuuuuuuuuuuuuuu",
				Longitude: point.Longitude(),
				Latitude:  point.Latitude(),
			},
			buildStus: func(repo *mock.MockUserRepository) {
				repo.EXPECT().
					SetUserLocation(
						gomock.Any(),
						gomock.Any(),
					).
					Times(0)
			},
			assert: func(t *testing.T, res domain.Location, err error) {
				require.Empty(t, res)
				var invalidArgumentErr *port.InvalidArgumentError
				require.ErrorAs(t, err, &invalidArgumentErr)
			},
		},
		{
			name: "Err_InvalidUsername_DoesNotMatchPattern",
			arg: port.SetUserLocationRequest{
				Username:  "user1_",
				Longitude: point.Longitude(),
				Latitude:  point.Latitude(),
			},
			buildStus: func(repo *mock.MockUserRepository) {
				repo.EXPECT().
					SetUserLocation(
						gomock.Any(),
						gomock.Any(),
					).
					Times(0)
			},
			assert: func(t *testing.T, res domain.Location, err error) {
				require.Empty(t, res)
				var invalidArgumentErr *port.InvalidArgumentError
				require.ErrorAs(t, err, &invalidArgumentErr)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.T().Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(s.T())
			defer ctrl.Finish()

			repo := mock.NewMockUserRepository(ctrl)
			tc.buildStus(repo)

			svc := service.NewUserService(repo)

			_, _ = svc.SetUserLocation(context.Background(), tc.arg)
		})
	}
}
