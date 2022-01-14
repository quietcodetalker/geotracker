package service_test

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gitlab.com/spacewalker/locations/internal/app/location/core/domain"
	"gitlab.com/spacewalker/locations/internal/app/location/core/port"
	"gitlab.com/spacewalker/locations/internal/app/location/core/port/mock"
	"gitlab.com/spacewalker/locations/internal/app/location/core/service"
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
		name       string
		arg        port.UserServiceSetUserLocationRequest
		buildStubs func(repo *mock.MockUserRepository)
		assert     func(t *testing.T, res domain.Location, err error)
	}{
		{
			name: "OK",
			arg: port.UserServiceSetUserLocationRequest{
				Username:  username,
				Longitude: point.Longitude(),
				Latitude:  point.Latitude(),
			},
			buildStubs: func(repo *mock.MockUserRepository) {
				repo.EXPECT().
					SetUserLocation(
						gomock.Any(),
						gomock.Eq(port.UserRepositorySetUserLocationRequest{
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
			arg: port.UserServiceSetUserLocationRequest{
				Username:  username,
				Longitude: -180.0,
				Latitude:  -90.0,
			},
			buildStubs: func(repo *mock.MockUserRepository) {
				repo.EXPECT().
					SetUserLocation(
						gomock.Any(),
						gomock.Eq(port.UserRepositorySetUserLocationRequest{
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
			arg: port.UserServiceSetUserLocationRequest{
				Username:  username,
				Longitude: 180.0,
				Latitude:  90.0,
			},
			buildStubs: func(repo *mock.MockUserRepository) {
				repo.EXPECT().
					SetUserLocation(
						gomock.Any(),
						gomock.Eq(port.UserRepositorySetUserLocationRequest{
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
			arg: port.UserServiceSetUserLocationRequest{
				Username:  username,
				Longitude: 0.0,
				Latitude:  0.0,
			},
			buildStubs: func(repo *mock.MockUserRepository) {
				repo.EXPECT().
					SetUserLocation(
						gomock.Any(),
						gomock.Eq(port.UserRepositorySetUserLocationRequest{
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
			arg: port.UserServiceSetUserLocationRequest{
				Username:  username,
				Longitude: -180.01,
				Latitude:  point.Latitude(),
			},
			buildStubs: func(repo *mock.MockUserRepository) {
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
			arg: port.UserServiceSetUserLocationRequest{
				Username:  username,
				Longitude: 180.01,
				Latitude:  point.Latitude(),
			},
			buildStubs: func(repo *mock.MockUserRepository) {
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
			arg: port.UserServiceSetUserLocationRequest{
				Username:  username,
				Longitude: point.Longitude(),
				Latitude:  -90.01,
			},
			buildStubs: func(repo *mock.MockUserRepository) {
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
			arg: port.UserServiceSetUserLocationRequest{
				Username:  username,
				Longitude: point.Longitude(),
				Latitude:  90.01,
			},
			buildStubs: func(repo *mock.MockUserRepository) {
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
			arg: port.UserServiceSetUserLocationRequest{
				Username:  "u",
				Longitude: point.Longitude(),
				Latitude:  point.Latitude(),
			},
			buildStubs: func(repo *mock.MockUserRepository) {
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
			arg: port.UserServiceSetUserLocationRequest{
				Username:  "uuuuuuuuuuuuuuuuu",
				Longitude: point.Longitude(),
				Latitude:  point.Latitude(),
			},
			buildStubs: func(repo *mock.MockUserRepository) {
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
			arg: port.UserServiceSetUserLocationRequest{
				Username:  "user1_",
				Longitude: point.Longitude(),
				Latitude:  point.Latitude(),
			},
			buildStubs: func(repo *mock.MockUserRepository) {
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
			tc.buildStubs(repo)

			svc := service.NewUserService(repo)

			_, _ = svc.SetUserLocation(context.Background(), tc.arg)
		})
	}
}

func (s *UserSvcTestSuite) Test_UserService_ListUsersInRadius() {
	testCases := []struct {
		name       string
		req        port.UserServiceListUsersInRadiusRequest
		buildStubs func(repo *mock.MockUserRepository)
		assert     func(t *testing.T, res port.UserServiceListUsersInRadiusResponse, err error)
	}{
		{
			name: "OK_PageToken",
			req: port.UserServiceListUsersInRadiusRequest{
				Point:     domain.Point{0, 0},
				Radius:    0,
				PageToken: "MTAwIDEwMA==",
				PageSize:  0,
			},
			buildStubs: func(repo *mock.MockUserRepository) {
				repo.EXPECT().
					ListUsersInRadius(gomock.Any(), gomock.Eq(port.UserRepositoryListUsersInRadiusRequest{
						Point:     domain.Point{0, 0},
						Radius:    0,
						PageToken: 100,
						PageSize:  100,
					})).
					Times(1).
					Return(port.UserRepositoryListUsersInRadiusResponse{
						Users: []domain.User{
							{
								ID:        1,
								Username:  "test",
								CreatedAt: time.Now(),
								UpdatedAt: time.Now(),
							},
						},
						NextPageToken: 1,
					}, nil)
			},
			assert: func(t *testing.T, res port.UserServiceListUsersInRadiusResponse, err error) {
				require.NoError(t, err)
				require.Len(t, res.Users, 1)
				require.Equal(t, 1, res.Users[0].ID)
				require.Equal(t, "test", res.Users[0].Username)
				require.WithinDuration(t, time.Now(), res.Users[0].CreatedAt, time.Second)
				require.WithinDuration(t, time.Now(), res.Users[0].UpdatedAt, time.Second)
				require.Equal(t, "MSAxMDA=", res.NextPageToken)
			},
		},
		{
			name: "OK_Lastpage",
			req: port.UserServiceListUsersInRadiusRequest{
				Point:     domain.Point{0, 0},
				Radius:    0,
				PageToken: "MTAwIDEwMA==",
				PageSize:  0,
			},
			buildStubs: func(repo *mock.MockUserRepository) {
				repo.EXPECT().
					ListUsersInRadius(gomock.Any(), gomock.Eq(port.UserRepositoryListUsersInRadiusRequest{
						Point:     domain.Point{0, 0},
						Radius:    0,
						PageToken: 100,
						PageSize:  100,
					})).
					Times(1).
					Return(port.UserRepositoryListUsersInRadiusResponse{
						Users: []domain.User{
							{
								ID:        1,
								Username:  "test",
								CreatedAt: time.Now(),
								UpdatedAt: time.Now(),
							},
						},
						NextPageToken: 0,
					}, nil)
			},
			assert: func(t *testing.T, res port.UserServiceListUsersInRadiusResponse, err error) {
				require.NoError(t, err)
				require.Len(t, res.Users, 1)
				require.Equal(t, 1, res.Users[0].ID)
				require.Equal(t, "test", res.Users[0].Username)
				require.WithinDuration(t, time.Now(), res.Users[0].CreatedAt, time.Second)
				require.WithinDuration(t, time.Now(), res.Users[0].UpdatedAt, time.Second)
				require.Equal(t, "", res.NextPageToken)
			},
		},
		{
			name: "OK_PageSize",
			req: port.UserServiceListUsersInRadiusRequest{
				Point:     domain.Point{0, 0},
				Radius:    0,
				PageToken: "",
				PageSize:  10,
			},
			buildStubs: func(repo *mock.MockUserRepository) {
				repo.EXPECT().
					ListUsersInRadius(gomock.Any(), gomock.Eq(port.UserRepositoryListUsersInRadiusRequest{
						Point:     domain.Point{0, 0},
						Radius:    0,
						PageToken: 0,
						PageSize:  10,
					})).
					Times(1).
					Return(port.UserRepositoryListUsersInRadiusResponse{
						Users: []domain.User{
							{
								ID:        1,
								Username:  "test",
								CreatedAt: time.Now(),
								UpdatedAt: time.Now(),
							},
						},
						NextPageToken: 1,
					}, nil)
			},
			assert: func(t *testing.T, res port.UserServiceListUsersInRadiusResponse, err error) {
				require.NoError(t, err)
				require.Len(t, res.Users, 1)
				require.Equal(t, 1, res.Users[0].ID)
				require.Equal(t, "test", res.Users[0].Username)
				require.WithinDuration(t, time.Now(), res.Users[0].CreatedAt, time.Second)
				require.WithinDuration(t, time.Now(), res.Users[0].UpdatedAt, time.Second)
				require.Equal(t, "MSAxMA==", res.NextPageToken)
			},
		},
		{
			name: "InvalidPoint_LongitudeLessThanMin",
			req: port.UserServiceListUsersInRadiusRequest{
				Point:     domain.Point{-180.1, 0},
				Radius:    0,
				PageToken: "1",
				PageSize:  0,
			},
			buildStubs: func(repo *mock.MockUserRepository) {
				repo.EXPECT().ListUsersInRadius(gomock.Any(), gomock.Any()).Times(0)
			},
			assert: func(t *testing.T, res port.UserServiceListUsersInRadiusResponse, err error) {
				require.Empty(t, res)

				var invalidArgumentError *port.InvalidArgumentError
				require.ErrorAs(t, err, &invalidArgumentError)
			},
		},
		{
			name: "InvalidPoint_LongitudeGreaterThanMax",
			req: port.UserServiceListUsersInRadiusRequest{
				Point:     domain.Point{180.1, 0},
				Radius:    0,
				PageToken: "1",
				PageSize:  0,
			},
			buildStubs: func(repo *mock.MockUserRepository) {
				repo.EXPECT().ListUsersInRadius(gomock.Any(), gomock.Any()).Times(0)
			},
			assert: func(t *testing.T, res port.UserServiceListUsersInRadiusResponse, err error) {
				require.Empty(t, res)

				var invalidArgumentError *port.InvalidArgumentError
				require.ErrorAs(t, err, &invalidArgumentError)
			},
		},
		{
			name: "InvalidPoint_LatitudeLessThanMin",
			req: port.UserServiceListUsersInRadiusRequest{
				Point:     domain.Point{0, -90.1},
				Radius:    0,
				PageToken: "1",
				PageSize:  0,
			},
			buildStubs: func(repo *mock.MockUserRepository) {
				repo.EXPECT().ListUsersInRadius(gomock.Any(), gomock.Any()).Times(0)
			},
			assert: func(t *testing.T, res port.UserServiceListUsersInRadiusResponse, err error) {
				require.Empty(t, res)

				var invalidArgumentError *port.InvalidArgumentError
				require.ErrorAs(t, err, &invalidArgumentError)
			},
		},
		{
			name: "InvalidPoint_LatitudeGreaterThanMax",
			req: port.UserServiceListUsersInRadiusRequest{
				Point:     domain.Point{0, 90.1},
				Radius:    0,
				PageToken: "1",
				PageSize:  0,
			},
			buildStubs: func(repo *mock.MockUserRepository) {
				repo.EXPECT().ListUsersInRadius(gomock.Any(), gomock.Any()).Times(0)
			},
			assert: func(t *testing.T, res port.UserServiceListUsersInRadiusResponse, err error) {
				require.Empty(t, res)

				var invalidArgumentError *port.InvalidArgumentError
				require.ErrorAs(t, err, &invalidArgumentError)
			},
		},
		{
			name: "PageTokenAndPageSizeBothProvided",
			req: port.UserServiceListUsersInRadiusRequest{
				Point:     domain.Point{0, 0},
				Radius:    0,
				PageToken: "1",
				PageSize:  10,
			},
			buildStubs: func(repo *mock.MockUserRepository) {
				repo.EXPECT().ListUsersInRadius(gomock.Any(), gomock.Any()).Times(0)
			},
			assert: func(t *testing.T, res port.UserServiceListUsersInRadiusResponse, err error) {
				require.Empty(t, res)

				var invalidArgumentError *port.InvalidArgumentError
				require.ErrorAs(t, err, &invalidArgumentError)
			},
		},
		{
			name: "PageTokenAndPageSizeBothNotProvided",
			req: port.UserServiceListUsersInRadiusRequest{
				Point:     domain.Point{0, 0},
				Radius:    0,
				PageToken: "1",
				PageSize:  10,
			},
			buildStubs: func(repo *mock.MockUserRepository) {
				repo.EXPECT().ListUsersInRadius(gomock.Any(), gomock.Any()).Times(0)
			},
			assert: func(t *testing.T, res port.UserServiceListUsersInRadiusResponse, err error) {
				require.Empty(t, res)

				var invalidArgumentError *port.InvalidArgumentError
				require.ErrorAs(t, err, &invalidArgumentError)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.T().Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mock.NewMockUserRepository(ctrl)
			tc.buildStubs(repo)

			svc := service.NewUserService(repo)

			res, err := svc.ListUsersInRadius(context.Background(), tc.req)

			tc.assert(t, res, err)
		})
	}
}
