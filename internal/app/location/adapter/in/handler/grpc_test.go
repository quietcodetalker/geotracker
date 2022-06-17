package handler_test

import (
  "context"
  "github.com/golang/mock/gomock"
  "github.com/stretchr/testify/require"
  "github.com/stretchr/testify/suite"
  "gitlab.com/spacewalker/geotracker/internal/app/location/adapter/in/handler"
  "gitlab.com/spacewalker/geotracker/internal/app/location/core/domain"
  "gitlab.com/spacewalker/geotracker/internal/app/location/core/port/mock"
  "gitlab.com/spacewalker/geotracker/internal/app/location/core/service"
  "gitlab.com/spacewalker/geotracker/internal/pkg/errpack"
  "gitlab.com/spacewalker/geotracker/internal/pkg/log"
  "gitlab.com/spacewalker/geotracker/internal/pkg/util/testutil"
  pb "gitlab.com/spacewalker/geotracker/pkg/api/proto/v1/location"
  "google.golang.org/grpc"
  "google.golang.org/grpc/codes"
  "google.golang.org/grpc/status"
  "google.golang.org/grpc/test/bufconn"
  "net"
  "testing"
  "time"
)

type GRPCHandlerTestSuite struct {
  suite.Suite
}

func (s *GRPCHandlerTestSuite) TestGetUserByUsername() {
  user := domain.User{
    ID:        testutil.RandomInt(1, 100),
    Username:  testutil.RandomUsername(),
    CreatedAt: time.Now(),
    UpdatedAt: time.Now(),
  }

  testCases := []struct {
    name            string
    buildStubs      func(svc *mock.MockUserRepository)
    req             *pb.GetUserByUsernameRequest
    expectedRes     *domain.User
    expectedErrCode codes.Code
  }{
    {
      name: "OK",
      buildStubs: func(svc *mock.MockUserRepository) {
        svc.EXPECT().
          GetByUsername(gomock.Any(), gomock.Eq(user.Username)).
          Times(1).
          Return(user, nil)
      },
      req:             &pb.GetUserByUsernameRequest{Username: user.Username},
      expectedRes:     &user,
      expectedErrCode: codes.OK,
    },
    {
      name: "invalid argument",
      buildStubs: func(svc *mock.MockUserRepository) {
        svc.EXPECT().
          GetByUsername(gomock.Any(), gomock.Any()).
          Times(0)
      },
      req:             &pb.GetUserByUsernameRequest{Username: ""},
      expectedErrCode: codes.InvalidArgument,
    },
    {
      name: "internal error",
      buildStubs: func(svc *mock.MockUserRepository) {
        svc.EXPECT().
          GetByUsername(gomock.Any(), gomock.Eq(user.Username)).
          Times(1).
          Return(domain.User{}, errpack.ErrInternalError)
      },
      req:             &pb.GetUserByUsernameRequest{Username: user.Username},
      expectedErrCode: codes.Internal,
    },
  }

  for _, tc := range testCases {
    tc := tc
    s.Run(tc.name, func() {
      ctrl := gomock.NewController(s.T())
      defer ctrl.Finish()

      repo := mock.NewMockUserRepository(ctrl)
      tc.buildStubs(repo)

      hc := mock.NewMockHistoryClient(ctrl)
      l := log.NewTestingLogger()

      svc := service.NewUserService(repo, hc, l)

      listener := bufconn.Listen(1024 * 1024)
      server := grpc.NewServer()
      pb.RegisterLocationInternalServer(server, handler.NewGRPCHandler(svc))

      go func() {
        if err := server.Serve(listener); err != nil {
          s.Fail(err.Error())
        }
      }()

      dial := func(context.Context, string) (net.Conn, error) {
        return listener.Dial()
      }

      ctx := context.Background()
      conn, err := grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(dial))
      if err != nil {
        s.Fail(err.Error())
      }
      defer conn.Close()

      client := pb.NewLocationInternalClient(conn)

      response, err := client.GetUserByUsername(context.Background(), tc.req)
      if tc.expectedErrCode == codes.OK {
        require.NoError(s.T(), err)

        require.Equal(s.T(), tc.expectedRes.ID, int(response.Id))
        require.Equal(s.T(), tc.expectedRes.Username, response.Username)

        require.WithinDuration(s.T(), tc.expectedRes.CreatedAt, response.CreatedAt.AsTime(), time.Second)
        require.WithinDuration(s.T(), tc.expectedRes.UpdatedAt, response.UpdatedAt.AsTime(), time.Second)
      }

      st, ok := status.FromError(err)
      require.True(s.T(), ok)
      st.Code()
    })
  }
}

func TestGRPCHandlerTestSuite(t *testing.T) {
  suite.Run(t, new(GRPCHandlerTestSuite))
}

func StartGRPCServer() {
}
