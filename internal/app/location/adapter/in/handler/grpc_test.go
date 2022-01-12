package handler_test

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"gitlab.com/spacewalker/locations/internal/app/location/adapter/in/handler"
	"gitlab.com/spacewalker/locations/internal/app/location/core/port"
	"gitlab.com/spacewalker/locations/internal/app/location/core/port/mock"
	pb "gitlab.com/spacewalker/locations/pkg/api/proto/v1/location"
	"testing"
)

func Test_GRPCHandler_SetUserLocation(t *testing.T) {
	testCases := []*struct {
		name       string
		req        *pb.SetUserLocationRequest
		buildStubs func(svc *mock.MockUserService)
		assert     func(t *testing.T, res *pb.SetUserLocationResponse, err error)
	}{
		{
			name: "OK",
			req: &pb.SetUserLocationRequest{
				Username:  "test",
				Longitude: 1.0,
				Latitude:  1.0,
			},
			buildStubs: func(svc *mock.MockUserService) {
				svc.EXPECT().
					SetUserLocation(
						gomock.Any(),
						gomock.Eq(port.UserServiceSetUserLocationRequest{
							Username:  "test",
							Latitude:  1.0,
							Longitude: 1.0,
						}),
					).
					Times(1).
					Return(port.UserServiceSetUserLocationResponse{
						Latitude:  1.0,
						Longitude: 1.0,
					}, nil)
			},
			assert: func(t *testing.T, res *pb.SetUserLocationResponse, err error) {
				require.NoError(t, err)
				require.Equal(t, 1.0, res.GetLongitude())
				require.Equal(t, 1.0, res.GetLatitude())
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			svc := mock.NewMockUserService(ctrl)
			tc.buildStubs(svc)

			hdl := handler.NewGRPCHandler(svc)

			res, err := hdl.SetUserLocation(context.Background(), tc.req)
			tc.assert(t, res, err)
		})
	}
}
