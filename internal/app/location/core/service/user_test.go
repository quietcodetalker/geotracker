package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gitlab.com/spacewalker/locations/internal/app/location/core/domain"
	"gitlab.com/spacewalker/locations/internal/app/location/core/port"
	"gitlab.com/spacewalker/locations/internal/app/location/core/port/mock"
	"gitlab.com/spacewalker/locations/internal/app/location/core/service"
	"gitlab.com/spacewalker/locations/internal/pkg/errpack"
	"gitlab.com/spacewalker/locations/internal/pkg/geo"
	mocklog "gitlab.com/spacewalker/locations/internal/pkg/log/mock"
	"gitlab.com/spacewalker/locations/internal/pkg/util"
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

	point := geo.Trunc(geo.Point{
		util.RandomFloat64(-180.0, 180.0),
		util.RandomFloat64(-90.0, 90.0),
	})

	testCases := []struct {
		name       string
		arg        port.UserServiceSetUserLocationRequest
		buildStubs func(repo *mock.MockUserRepository, historyClient *mock.MockHistoryClient)
		assert     func(t *testing.T, res domain.Location, err error)
	}{
		{
			name: "OK",
			arg: port.UserServiceSetUserLocationRequest{
				Username:  username,
				Longitude: point.Longitude(),
				Latitude:  point.Latitude(),
			},
			buildStubs: func(repo *mock.MockUserRepository, historyClient *mock.MockHistoryClient) {
				historyClient.EXPECT().AddRecord(gomock.Any(), gomock.Any())
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
						port.UserRepositorySetUserLocationResponse{
							Location: domain.Location{
								UserID:    1,
								Point:     point,
								CreatedAt: time.Now(),
								UpdatedAt: time.Now(),
							},
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
			buildStubs: func(repo *mock.MockUserRepository, historyClient *mock.MockHistoryClient) {
				historyClient.EXPECT().AddRecord(gomock.Any(), gomock.Any())
				repo.EXPECT().
					SetUserLocation(
						gomock.Any(),
						gomock.Eq(port.UserRepositorySetUserLocationRequest{
							Username: username,
							Point:    geo.Point{-180.0, -90.0},
						}),
					).
					Times(1).
					Return(
						port.UserRepositorySetUserLocationResponse{
							Location: domain.Location{
								UserID:    1,
								Point:     geo.Point{-180.0, -90.0},
								CreatedAt: time.Now(),
								UpdatedAt: time.Now(),
							},
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
			buildStubs: func(repo *mock.MockUserRepository, historyClient *mock.MockHistoryClient) {
				historyClient.EXPECT().AddRecord(gomock.Any(), gomock.Any())
				repo.EXPECT().
					SetUserLocation(
						gomock.Any(),
						gomock.Eq(port.UserRepositorySetUserLocationRequest{
							Username: username,
							Point:    geo.Point{180.0, 90.0},
						}),
					).
					Times(1).
					Return(
						port.UserRepositorySetUserLocationResponse{
							Location: domain.Location{
								UserID:    1,
								Point:     geo.Point{180.0, 90.0},
								CreatedAt: time.Now(),
								UpdatedAt: time.Now(),
							},
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
			buildStubs: func(repo *mock.MockUserRepository, historyClient *mock.MockHistoryClient) {
				historyClient.EXPECT().AddRecord(gomock.Any(), gomock.Any())
				repo.EXPECT().
					SetUserLocation(
						gomock.Any(),
						gomock.Eq(port.UserRepositorySetUserLocationRequest{
							Username: username,
							Point:    geo.Point{0, 0},
						}),
					).
					Times(1).
					Return(
						port.UserRepositorySetUserLocationResponse{
							Location: domain.Location{
								UserID:    1,
								Point:     geo.Point{0, 0},
								CreatedAt: time.Now(),
								UpdatedAt: time.Now(),
							},
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
			buildStubs: func(repo *mock.MockUserRepository, historyClient *mock.MockHistoryClient) {
				repo.EXPECT().
					SetUserLocation(
						gomock.Any(),
						gomock.Any(),
					).
					Times(0)
			},
			assert: func(t *testing.T, res domain.Location, err error) {
				require.Empty(t, res)
				require.ErrorIs(t, err, errpack.ErrInvalidArgument)
			},
		},
		{
			name: "Err_InvalidLongitude_GreaterThanMax",
			arg: port.UserServiceSetUserLocationRequest{
				Username:  username,
				Longitude: 180.01,
				Latitude:  point.Latitude(),
			},
			buildStubs: func(repo *mock.MockUserRepository, historyClient *mock.MockHistoryClient) {
				repo.EXPECT().
					SetUserLocation(
						gomock.Any(),
						gomock.Any(),
					).
					Times(0)
			},
			assert: func(t *testing.T, res domain.Location, err error) {
				require.Empty(t, res)
				require.ErrorIs(t, err, errpack.ErrInvalidArgument)
			},
		},
		{
			name: "Err_InvalidLatitude_LessThanMin",
			arg: port.UserServiceSetUserLocationRequest{
				Username:  username,
				Longitude: point.Longitude(),
				Latitude:  -90.01,
			},
			buildStubs: func(repo *mock.MockUserRepository, historyClient *mock.MockHistoryClient) {
				repo.EXPECT().
					SetUserLocation(
						gomock.Any(),
						gomock.Any(),
					).
					Times(0)
			},
			assert: func(t *testing.T, res domain.Location, err error) {
				require.Empty(t, res)
				require.ErrorIs(t, err, errpack.ErrInvalidArgument)
			},
		},
		{
			name: "Err_InvalidLatitude_GreaterThanMax",
			arg: port.UserServiceSetUserLocationRequest{
				Username:  username,
				Longitude: point.Longitude(),
				Latitude:  90.01,
			},
			buildStubs: func(repo *mock.MockUserRepository, historyClient *mock.MockHistoryClient) {
				repo.EXPECT().
					SetUserLocation(
						gomock.Any(),
						gomock.Any(),
					).
					Times(0)
			},
			assert: func(t *testing.T, res domain.Location, err error) {
				require.Empty(t, res)
				require.ErrorIs(t, err, errpack.ErrInvalidArgument)
			},
		},
		{
			name: "Err_InvalidUsername_TooShort",
			arg: port.UserServiceSetUserLocationRequest{
				Username:  "u",
				Longitude: point.Longitude(),
				Latitude:  point.Latitude(),
			},
			buildStubs: func(repo *mock.MockUserRepository, historyClient *mock.MockHistoryClient) {
				repo.EXPECT().
					SetUserLocation(
						gomock.Any(),
						gomock.Any(),
					).
					Times(0)
			},
			assert: func(t *testing.T, res domain.Location, err error) {
				require.Empty(t, res)
				require.ErrorIs(t, err, errpack.ErrInvalidArgument)
			},
		},
		{
			name: "Err_InvalidUsername_TooLong",
			arg: port.UserServiceSetUserLocationRequest{
				Username:  "uuuuuuuuuuuuuuuuu",
				Longitude: point.Longitude(),
				Latitude:  point.Latitude(),
			},
			buildStubs: func(repo *mock.MockUserRepository, historyClient *mock.MockHistoryClient) {
				repo.EXPECT().
					SetUserLocation(
						gomock.Any(),
						gomock.Any(),
					).
					Times(0)
			},
			assert: func(t *testing.T, res domain.Location, err error) {
				require.Empty(t, res)
				require.ErrorIs(t, err, errpack.ErrInvalidArgument)
			},
		},
		{
			name: "Err_InvalidUsername_DoesNotMatchPattern",
			arg: port.UserServiceSetUserLocationRequest{
				Username:  "user1_",
				Longitude: point.Longitude(),
				Latitude:  point.Latitude(),
			},
			buildStubs: func(repo *mock.MockUserRepository, historyClient *mock.MockHistoryClient) {
				repo.EXPECT().
					SetUserLocation(
						gomock.Any(),
						gomock.Any(),
					).
					Times(0)
			},
			assert: func(t *testing.T, res domain.Location, err error) {
				require.Empty(t, res)
				require.ErrorIs(t, err, errpack.ErrInvalidArgument)
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.T().Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(s.T())
			defer ctrl.Finish()

			repo := mock.NewMockUserRepository(ctrl)
			historyClient := mock.NewMockHistoryClient(ctrl)
			tc.buildStubs(repo, historyClient)
			logger := mocklog.NewMockLogger(ctrl)
			svc := service.NewUserService(repo, historyClient, logger)

			_, _ = svc.SetUserLocation(context.Background(), tc.arg)
		})
	}
}

func (s *UserSvcTestSuite) Test_UserService_ListUsersInRadius() {
	testCases := []struct {
		name       string
		req        port.UserServiceListUsersInRadiusRequest
		buildStubs func(repo *mock.MockUserRepository, historyClient *mock.MockHistoryClient)
		assert     func(t *testing.T, res port.UserServiceListUsersInRadiusResponse, err error)
	}{
		{
			name: "OK_PageToken",
			req: port.UserServiceListUsersInRadiusRequest{
				Point:     geo.Point{0, 0},
				Radius:    0,
				PageToken: "MTAwIDEwMA==",
				PageSize:  0,
			},
			buildStubs: func(repo *mock.MockUserRepository, historyClient *mock.MockHistoryClient) {
				repo.EXPECT().
					ListUsersInRadius(gomock.Any(), gomock.Eq(port.UserRepositoryListUsersInRadiusRequest{
						Point:     geo.Point{0, 0},
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
				Point:     geo.Point{0, 0},
				Radius:    0,
				PageToken: "MTAwIDEwMA==",
				PageSize:  0,
			},
			buildStubs: func(repo *mock.MockUserRepository, historyClient *mock.MockHistoryClient) {
				repo.EXPECT().
					ListUsersInRadius(gomock.Any(), gomock.Eq(port.UserRepositoryListUsersInRadiusRequest{
						Point:     geo.Point{0, 0},
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
				Point:     geo.Point{0, 0},
				Radius:    0,
				PageToken: "",
				PageSize:  10,
			},
			buildStubs: func(repo *mock.MockUserRepository, historyClient *mock.MockHistoryClient) {
				repo.EXPECT().
					ListUsersInRadius(gomock.Any(), gomock.Eq(port.UserRepositoryListUsersInRadiusRequest{
						Point:     geo.Point{0, 0},
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
				Point:     geo.Point{-180.1, 0},
				Radius:    0,
				PageToken: "1",
				PageSize:  0,
			},
			buildStubs: func(repo *mock.MockUserRepository, historyClient *mock.MockHistoryClient) {
				repo.EXPECT().ListUsersInRadius(gomock.Any(), gomock.Any()).Times(0)
			},
			assert: func(t *testing.T, res port.UserServiceListUsersInRadiusResponse, err error) {
				require.Empty(t, res)
				require.ErrorIs(t, err, errpack.ErrInvalidArgument)
			},
		},
		{
			name: "InvalidPoint_LongitudeGreaterThanMax",
			req: port.UserServiceListUsersInRadiusRequest{
				Point:     geo.Point{180.1, 0},
				Radius:    0,
				PageToken: "1",
				PageSize:  0,
			},
			buildStubs: func(repo *mock.MockUserRepository, historyClient *mock.MockHistoryClient) {
				repo.EXPECT().ListUsersInRadius(gomock.Any(), gomock.Any()).Times(0)
			},
			assert: func(t *testing.T, res port.UserServiceListUsersInRadiusResponse, err error) {
				require.Empty(t, res)
				require.ErrorIs(t, err, errpack.ErrInvalidArgument)
			},
		},
		{
			name: "InvalidPoint_LatitudeLessThanMin",
			req: port.UserServiceListUsersInRadiusRequest{
				Point:     geo.Point{0, -90.1},
				Radius:    0,
				PageToken: "1",
				PageSize:  0,
			},
			buildStubs: func(repo *mock.MockUserRepository, historyClient *mock.MockHistoryClient) {
				repo.EXPECT().ListUsersInRadius(gomock.Any(), gomock.Any()).Times(0)
			},
			assert: func(t *testing.T, res port.UserServiceListUsersInRadiusResponse, err error) {
				require.Empty(t, res)
				require.ErrorIs(t, err, errpack.ErrInvalidArgument)
			},
		},
		{
			name: "InvalidPoint_LatitudeGreaterThanMax",
			req: port.UserServiceListUsersInRadiusRequest{
				Point:     geo.Point{0, 90.1},
				Radius:    0,
				PageToken: "1",
				PageSize:  0,
			},
			buildStubs: func(repo *mock.MockUserRepository, historyClient *mock.MockHistoryClient) {
				repo.EXPECT().ListUsersInRadius(gomock.Any(), gomock.Any()).Times(0)
			},
			assert: func(t *testing.T, res port.UserServiceListUsersInRadiusResponse, err error) {
				require.Empty(t, res)
				require.ErrorIs(t, err, errpack.ErrInvalidArgument)
			},
		},
		{
			name: "PageTokenAndPageSizeBothProvided",
			req: port.UserServiceListUsersInRadiusRequest{
				Point:     geo.Point{0, 0},
				Radius:    0,
				PageToken: "1",
				PageSize:  10,
			},
			buildStubs: func(repo *mock.MockUserRepository, historyClient *mock.MockHistoryClient) {
				repo.EXPECT().ListUsersInRadius(gomock.Any(), gomock.Any()).Times(0)
			},
			assert: func(t *testing.T, res port.UserServiceListUsersInRadiusResponse, err error) {
				require.Empty(t, res)
				require.ErrorIs(t, err, errpack.ErrInvalidArgument)
			},
		},
		{
			name: "PageTokenAndPageSizeBothNotProvided",
			req: port.UserServiceListUsersInRadiusRequest{
				Point:     geo.Point{0, 0},
				Radius:    0,
				PageToken: "1",
				PageSize:  10,
			},
			buildStubs: func(repo *mock.MockUserRepository, historyClient *mock.MockHistoryClient) {
				repo.EXPECT().ListUsersInRadius(gomock.Any(), gomock.Any()).Times(0)
			},
			assert: func(t *testing.T, res port.UserServiceListUsersInRadiusResponse, err error) {
				require.Empty(t, res)
				require.ErrorIs(t, err, errpack.ErrInvalidArgument)
			},
		},
		{
			name: "InvalidRadius",
			req: port.UserServiceListUsersInRadiusRequest{
				Point:     geo.Point{0, 0},
				Radius:    -1,
				PageToken: "",
				PageSize:  10,
			},
			buildStubs: func(repo *mock.MockUserRepository, historyClient *mock.MockHistoryClient) {
				repo.EXPECT().ListUsersInRadius(gomock.Any(), gomock.Any()).Times(0)
			},
			assert: func(t *testing.T, res port.UserServiceListUsersInRadiusResponse, err error) {
				require.Empty(t, res)
				require.ErrorIs(t, err, errpack.ErrInvalidArgument)
			},
		},
		{
			name: "OK_Radius_eq_0",
			req: port.UserServiceListUsersInRadiusRequest{
				Point:     geo.Point{0, 0},
				Radius:    0,
				PageToken: "MTAwIDEwMA==",
				PageSize:  0,
			},
			buildStubs: func(repo *mock.MockUserRepository, historyClient *mock.MockHistoryClient) {
				repo.EXPECT().
					ListUsersInRadius(gomock.Any(), gomock.Eq(port.UserRepositoryListUsersInRadiusRequest{
						Point:     geo.Point{0, 0},
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
			name: "OK_Radius_gt_0",
			req: port.UserServiceListUsersInRadiusRequest{
				Point:     geo.Point{0, 0},
				Radius:    1,
				PageToken: "MTAwIDEwMA==",
				PageSize:  0,
			},
			buildStubs: func(repo *mock.MockUserRepository, historyClient *mock.MockHistoryClient) {
				repo.EXPECT().
					ListUsersInRadius(gomock.Any(), gomock.Eq(port.UserRepositoryListUsersInRadiusRequest{
						Point:     geo.Point{0, 0},
						Radius:    1,
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
	}

	for _, tc := range testCases {
		tc := tc
		s.T().Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mock.NewMockUserRepository(ctrl)
			historyClient := mock.NewMockHistoryClient(ctrl)
			tc.buildStubs(repo, historyClient)
			logger := mocklog.NewMockLogger(ctrl)
			svc := service.NewUserService(repo, historyClient, logger)

			res, err := svc.ListUsersInRadius(context.Background(), tc.req)

			tc.assert(t, res, err)
		})
	}
}

func (s *UserSvcTestSuite) Test_UserService_GetByUsername() {
	users := []domain.User{
		{
			ID:        1,
			Username:  "user1",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	errInternal := errors.New("internal error")

	testCases := []struct {
		name       string
		buildStubs func(repository *mock.MockUserRepository)
		username   string
		expected   domain.User
		hasError   bool
		isError    error
		asError    error
	}{
		{
			name: "OK",
			buildStubs: func(repository *mock.MockUserRepository) {
				repository.EXPECT().
					GetByUsername(gomock.Any(), gomock.Eq(users[0].Username)).
					Times(1).
					Return(users[0], nil)
			},
			username: users[0].Username,
			expected: users[0],
			hasError: false,
		},
		{
			name: "InvalidUsername",
			buildStubs: func(repository *mock.MockUserRepository) {
				repository.EXPECT().
					GetByUsername(gomock.Any(), gomock.Any()).
					Times(0)
			},
			username: "",
			expected: domain.User{},
			hasError: true,
			isError:  errpack.ErrInvalidArgument,
		},
		{
			name: "UserNotFound",
			buildStubs: func(repository *mock.MockUserRepository) {
				repository.EXPECT().
					GetByUsername(gomock.Any(), gomock.Eq(users[0].Username)).
					Times(1).
					Return(domain.User{}, errpack.ErrNotFound)
			},
			username: users[0].Username,
			expected: domain.User{},
			hasError: true,
			isError:  errpack.ErrNotFound,
		},
		{
			name: "InternalError",
			buildStubs: func(repository *mock.MockUserRepository) {
				repository.EXPECT().
					GetByUsername(gomock.Any(), gomock.Eq(users[0].Username)).
					Times(1).
					Return(domain.User{}, errInternal)
			},
			username: users[0].Username,
			expected: domain.User{},
			hasError: true,
			isError:  errInternal,
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.T().Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mock.NewMockUserRepository(ctrl)
			historyClient := mock.NewMockHistoryClient(ctrl)
			tc.buildStubs(repo)
			logger := mocklog.NewMockLogger(ctrl)
			svc := service.NewUserService(repo, historyClient, logger)

			user, err := svc.GetByUsername(context.Background(), tc.username)
			if tc.hasError {
				if tc.isError != nil {
					require.ErrorIs(t, err, tc.isError)
				}
				if tc.asError != nil {
					require.ErrorAs(t, err, &tc.asError)
				}
			} else {
				require.NoError(t, err)

				require.Equal(t, tc.expected.ID, user.ID)
				require.Equal(t, tc.expected.Username, user.Username)
				require.Equal(t, tc.expected.CreatedAt, user.CreatedAt)
				require.Equal(t, tc.expected.UpdatedAt, user.UpdatedAt)
			}
		})
	}
}
