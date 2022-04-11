package handler_test

import (
  "context"
  "fmt"
  "github.com/golang/mock/gomock"
  "github.com/stretchr/testify/require"
  "github.com/stretchr/testify/suite"
  "gitlab.com/spacewalker/geotracker/internal/app/history/adapter/in/handler"
  "gitlab.com/spacewalker/geotracker/internal/app/history/core/domain"
  "gitlab.com/spacewalker/geotracker/internal/app/history/core/port"
  "gitlab.com/spacewalker/geotracker/internal/app/history/core/port/mock"
  "gitlab.com/spacewalker/geotracker/internal/app/history/core/service"
  "gitlab.com/spacewalker/geotracker/internal/pkg/errpack"
  "gitlab.com/spacewalker/geotracker/internal/pkg/geo"
  "gitlab.com/spacewalker/geotracker/internal/pkg/log"
  "gitlab.com/spacewalker/geotracker/internal/pkg/util/testutil"
  pb "gitlab.com/spacewalker/geotracker/pkg/api/proto/v1/history"
  "google.golang.org/grpc"
  "google.golang.org/grpc/codes"
  "google.golang.org/grpc/status"
  "google.golang.org/grpc/test/bufconn"
  "google.golang.org/protobuf/types/known/timestamppb"
  "net"
  "testing"
  "time"
)

type GRPCHandlerTestSuite struct {
  suite.Suite
}

func (s *GRPCHandlerTestSuite) TestAddRecord() {
  userID := testutil.RandomInt(1, 100)
  recordID := testutil.RandomInt(1, 100)
  a := geo.Point{
    testutil.RandomLongitude(),
    testutil.RandomLatitude(),
  }
  truncedA := geo.Trunc(a)
  b := geo.Point{
    testutil.RandomLongitude(),
    testutil.RandomLatitude(),
  }
  truncedB := geo.Trunc(b)
  timestamp := time.Now().UTC()

  testCases := []struct {
    name            string
    buildStubs      func(repo *mock.MockHistoryRepository)
    req             *pb.AddRecordRequest
    expectedRes     *pb.AddRecordResponse
    expectedErrCode codes.Code
  }{
    {
      name: "OK",
      buildStubs: func(repo *mock.MockHistoryRepository) {
        repo.EXPECT().
          AddRecord(gomock.Any(), gomock.Eq(port.HistoryRepositoryAddRecordRequest{
            UserID:    userID,
            A:         truncedA,
            B:         truncedB,
            Timestamp: timestamp,
          })).
          Times(1).
          Return(domain.Record{
            ID:        recordID,
            UserID:    userID,
            A:         truncedA,
            B:         truncedB,
            Timestamp: timestamp,
          }, nil)
      },
      req: &pb.AddRecordRequest{
        UserId: int32(userID),
        A: &pb.Point{
          Longitude: a.Longitude(),
          Latitude:  a.Latitude(),
        },
        B: &pb.Point{
          Longitude: b.Longitude(),
          Latitude:  b.Latitude(),
        },
        Timestamp: timestamppb.New(timestamp),
      },
      expectedRes: &pb.AddRecordResponse{
        UserId: int32(userID),
        A: &pb.Point{
          Longitude: truncedA.Longitude(),
          Latitude:  truncedA.Latitude(),
        },
        B: &pb.Point{
          Longitude: truncedB.Longitude(),
          Latitude:  truncedB.Latitude(),
        },
        Timestamp: timestamppb.New(timestamp),
      },
      expectedErrCode: codes.OK,
    },
    {
      name: "user not found",
      buildStubs: func(repo *mock.MockHistoryRepository) {
        repo.EXPECT().
          AddRecord(gomock.Any(), gomock.Eq(port.HistoryRepositoryAddRecordRequest{
            UserID:    userID,
            A:         truncedA,
            B:         truncedB,
            Timestamp: timestamp,
          })).
          Times(1).
          Return(domain.Record{}, fmt.Errorf("%w", errpack.ErrNotFound))
      },
      req: &pb.AddRecordRequest{
        UserId: int32(userID),
        A: &pb.Point{
          Longitude: a.Longitude(),
          Latitude:  a.Latitude(),
        },
        B: &pb.Point{
          Longitude: b.Longitude(),
          Latitude:  b.Latitude(),
        },
        Timestamp: timestamppb.New(timestamp),
      },
      expectedRes:     nil,
      expectedErrCode: codes.NotFound,
    },
    {
      name: "internal error",
      buildStubs: func(repo *mock.MockHistoryRepository) {
        repo.EXPECT().
          AddRecord(gomock.Any(), gomock.Eq(port.HistoryRepositoryAddRecordRequest{
            UserID:    userID,
            A:         truncedA,
            B:         truncedB,
            Timestamp: timestamp,
          })).
          Times(1).
          Return(domain.Record{}, fmt.Errorf("%w", errpack.ErrInternalError))
      },
      req: &pb.AddRecordRequest{
        UserId: int32(userID),
        A: &pb.Point{
          Longitude: a.Longitude(),
          Latitude:  a.Latitude(),
        },
        B: &pb.Point{
          Longitude: b.Longitude(),
          Latitude:  b.Latitude(),
        },
        Timestamp: timestamppb.New(timestamp),
      },
      expectedRes:     nil,
      expectedErrCode: codes.Internal,
    },
    {
      name: "invalid argument a long lt min",
      buildStubs: func(repo *mock.MockHistoryRepository) {
        repo.EXPECT().
          AddRecord(gomock.Any(), gomock.Any()).
          Times(0)
      },
      req: &pb.AddRecordRequest{
        UserId: int32(userID),
        A: &pb.Point{
          Longitude: -180.1,
          Latitude:  a.Latitude(),
        },
        B: &pb.Point{
          Longitude: b.Longitude(),
          Latitude:  b.Latitude(),
        },
        Timestamp: timestamppb.New(timestamp),
      },
      expectedRes:     nil,
      expectedErrCode: codes.InvalidArgument,
    },
    {
      name: "invalid argument a long gt max",
      buildStubs: func(repo *mock.MockHistoryRepository) {
        repo.EXPECT().
          AddRecord(gomock.Any(), gomock.Any()).
          Times(0)
      },
      req: &pb.AddRecordRequest{
        UserId: int32(userID),
        A: &pb.Point{
          Longitude: 180.1,
          Latitude:  a.Latitude(),
        },
        B: &pb.Point{
          Longitude: b.Longitude(),
          Latitude:  b.Latitude(),
        },
        Timestamp: timestamppb.New(timestamp),
      },
      expectedRes:     nil,
      expectedErrCode: codes.InvalidArgument,
    },
    {
      name: "invalid argument a lat lt min",
      buildStubs: func(repo *mock.MockHistoryRepository) {
        repo.EXPECT().
          AddRecord(gomock.Any(), gomock.Any()).
          Times(0)
      },
      req: &pb.AddRecordRequest{
        UserId: int32(userID),
        A: &pb.Point{
          Longitude: a.Longitude(),
          Latitude:  -90.1,
        },
        B: &pb.Point{
          Longitude: b.Longitude(),
          Latitude:  b.Latitude(),
        },
        Timestamp: timestamppb.New(timestamp),
      },
      expectedRes:     nil,
      expectedErrCode: codes.InvalidArgument,
    },
    {
      name: "invalid argument a lat gt max",
      buildStubs: func(repo *mock.MockHistoryRepository) {
        repo.EXPECT().
          AddRecord(gomock.Any(), gomock.Any()).
          Times(0)
      },
      req: &pb.AddRecordRequest{
        UserId: int32(userID),
        A: &pb.Point{
          Longitude: a.Longitude(),
          Latitude:  90.1,
        },
        B: &pb.Point{
          Longitude: b.Longitude(),
          Latitude:  b.Latitude(),
        },
        Timestamp: timestamppb.New(timestamp),
      },
      expectedRes:     nil,
      expectedErrCode: codes.InvalidArgument,
    },
    {
      name: "invalid argument b long lt min",
      buildStubs: func(repo *mock.MockHistoryRepository) {
        repo.EXPECT().
          AddRecord(gomock.Any(), gomock.Any()).
          Times(0)
      },
      req: &pb.AddRecordRequest{
        UserId: int32(userID),
        A: &pb.Point{
          Longitude: a.Longitude(),
          Latitude:  a.Latitude(),
        },
        B: &pb.Point{
          Longitude: -180.1,
          Latitude:  b.Latitude(),
        },
        Timestamp: timestamppb.New(timestamp),
      },
      expectedRes:     nil,
      expectedErrCode: codes.InvalidArgument,
    },
    {
      name: "invalid argument b long gt max",
      buildStubs: func(repo *mock.MockHistoryRepository) {
        repo.EXPECT().
          AddRecord(gomock.Any(), gomock.Any()).
          Times(0)
      },
      req: &pb.AddRecordRequest{
        UserId: int32(userID),
        A: &pb.Point{
          Longitude: a.Longitude(),
          Latitude:  a.Latitude(),
        },
        B: &pb.Point{
          Longitude: 180.1,
          Latitude:  b.Latitude(),
        },
        Timestamp: timestamppb.New(timestamp),
      },
      expectedRes:     nil,
      expectedErrCode: codes.InvalidArgument,
    },
    {
      name: "invalid argument b lat lt min",
      buildStubs: func(repo *mock.MockHistoryRepository) {
        repo.EXPECT().
          AddRecord(gomock.Any(), gomock.Any()).
          Times(0)
      },
      req: &pb.AddRecordRequest{
        UserId: int32(userID),
        A: &pb.Point{
          Longitude: a.Longitude(),
          Latitude:  a.Latitude(),
        },
        B: &pb.Point{
          Longitude: b.Longitude(),
          Latitude:  -90.1,
        },
        Timestamp: timestamppb.New(timestamp),
      },
      expectedRes:     nil,
      expectedErrCode: codes.InvalidArgument,
    },
    {
      name: "invalid argument b lat gt max",
      buildStubs: func(repo *mock.MockHistoryRepository) {
        repo.EXPECT().
          AddRecord(gomock.Any(), gomock.Any()).
          Times(0)
      },
      req: &pb.AddRecordRequest{
        UserId: int32(userID),
        A: &pb.Point{
          Longitude: a.Longitude(),
          Latitude:  a.Latitude(),
        },
        B: &pb.Point{
          Longitude: b.Longitude(),
          Latitude:  90.1,
        },
        Timestamp: timestamppb.New(timestamp),
      },
      expectedRes:     nil,
      expectedErrCode: codes.InvalidArgument,
    },
  }

  for _, tc := range testCases {
    tc := tc
    s.Run(tc.name, func() {
      ctrl := gomock.NewController(s.T())
      defer ctrl.Finish()

      repo := mock.NewMockHistoryRepository(ctrl)
      tc.buildStubs(repo)

      hc := mock.NewMockLocationClient(ctrl)
      l := log.NewTestingLogger()

      svc := service.NewHistoryService(repo, hc, l)

      listener := bufconn.Listen(1024 * 1024)
      server := grpc.NewServer()
      pb.RegisterHistoryServer(server, handler.NewGRPCHandler(svc))

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

      client := pb.NewHistoryClient(conn)

      response, err := client.AddRecord(context.Background(), tc.req)
      if tc.expectedErrCode == codes.OK {
        require.NoError(s.T(), err)

        require.Equal(s.T(), tc.expectedRes.UserId, response.UserId)
        require.Equal(s.T(), tc.expectedRes.A.Longitude, response.A.Longitude)
        require.Equal(s.T(), tc.expectedRes.A.Latitude, response.A.Latitude)
        require.Equal(s.T(), tc.expectedRes.B.Longitude, response.B.Longitude)
        require.Equal(s.T(), tc.expectedRes.B.Latitude, response.B.Latitude)

        require.WithinDuration(s.T(), tc.expectedRes.Timestamp.AsTime(), response.Timestamp.AsTime(), time.Second)
      }

      st, ok := status.FromError(err)
      require.True(s.T(), ok)
      st.Code()
    })
  }
}

func (s *GRPCHandlerTestSuite) TestGetDistance() {
  userID := testutil.RandomInt(1, 100)
  from, to := testutil.RandomTimeInterval()
  distance := testutil.RandomFloat64(0, 1000.0)

  testCases := []struct {
    name            string
    buildStubs      func(repo *mock.MockHistoryRepository)
    req             *pb.GetDistanceRequest
    expectedRes     *pb.GetDistanceResponse
    expectedErrCode codes.Code
  }{
    {
      name: "OK",
      buildStubs: func(repo *mock.MockHistoryRepository) {
        repo.EXPECT().
          GetDistance(gomock.Any(), port.HistoryRepositoryGetDistanceRequest{
            UserID: userID,
            From:   from,
            To:     to,
          }).
          Times(0)
      },
      req: &pb.GetDistanceRequest{
        UserId: int32(userID),
        From:   timestamppb.New(from),
        To:     timestamppb.New(to),
      },
      expectedRes: &pb.GetDistanceResponse{
        Distance: distance,
      },
      expectedErrCode: codes.OK,
    },
  }

  for _, tc := range testCases {
    tc := tc
    s.Run(tc.name, func() {
      ctrl := gomock.NewController(s.T())
      //defer ctrl.Finish()

      repo := mock.NewMockHistoryRepository(ctrl)
      tc.buildStubs(repo)

      hc := mock.NewMockLocationClient(ctrl)
      l := log.NewTestingLogger()

      svc := service.NewHistoryService(repo, hc, l)

      listener, err := net.Listen("tcp", ":")
      if err != nil {
        s.T().Fatalf("failed to listen: %v", err)
      }

      //listener := bufconn.Listen(1024 * 1024)
      server := grpc.NewServer()
      hndl := handler.NewGRPCHandler(svc)
      pb.RegisterHistoryServer(server, hndl)

      go func() {
        if err := server.Serve(listener); err != nil {
          s.T().Fatal()
        }
      }()

      //dial := func(context.Context, string) (net.Conn, error) {
      //  return listener.Dial()
      //}

      ctx, cancelCtx := context.WithTimeout(context.Background(), time.Second)
      defer cancelCtx()
      conn, err := grpc.DialContext(ctx, listener.Addr().String(), grpc.WithInsecure())
      if err != nil {
        s.T().Fatal()
      }
      defer conn.Close()

      client := pb.NewHistoryClient(conn)

      response, err := client.GetDistance(context.Background(), tc.req)
      if tc.expectedErrCode == codes.OK {
        require.NoError(s.T(), err)
        require.Equal(s.T(), tc.expectedRes.Distance, response.Distance)
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
