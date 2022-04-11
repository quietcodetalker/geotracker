package handler_test

import (
  "bytes"
  "encoding/json"
  "fmt"
  "github.com/gavv/httpexpect"
  "github.com/golang/mock/gomock"
  "github.com/stretchr/testify/suite"
  "gitlab.com/spacewalker/geotracker/internal/app/location/adapter/in/handler"
  "gitlab.com/spacewalker/geotracker/internal/app/location/core/domain"
  "gitlab.com/spacewalker/geotracker/internal/app/location/core/port"
  "gitlab.com/spacewalker/geotracker/internal/app/location/core/port/mock"
  "gitlab.com/spacewalker/geotracker/internal/app/location/core/service"
  "gitlab.com/spacewalker/geotracker/internal/pkg/errpack"
  "gitlab.com/spacewalker/geotracker/internal/pkg/geo"
  "gitlab.com/spacewalker/geotracker/internal/pkg/log"
  "gitlab.com/spacewalker/geotracker/internal/pkg/util/testutil"
  "net/http"
  "net/http/httptest"
  "testing"
)

type HTTPHandleTestSuite struct {
  suite.Suite
}

func (s *HTTPHandleTestSuite) TestSetUserLocation() {
  path := "/users/{username}/location"
  username := testutil.RandomUsername()
  shortUsername := testutil.RandomString(
    1,
    testutil.UsernameMinLen-1,
    testutil.UsernameCharacterSet,
  )
  longUsername := testutil.RandomString(
    testutil.UsernameMaxLen+1,
    testutil.UsernameMaxLen+2,
    testutil.UsernameCharacterSet,
  )
  invalidUsername := testutil.RandomString(
    testutil.UsernameMinLen,
    testutil.UsernameMaxLen,
    "!@##%&*^",
  )
  latitude := testutil.RandomLatitude()
  longitude := testutil.RandomLongitude()

  buildStubsNoCallExpected := func(repo *mock.MockUserRepository) {
    repo.EXPECT().
      SetUserLocation(gomock.Any(), gomock.Any()).
      Times(0)
  }

  invalidArguentResponse := map[string]interface{}{
    "error": map[string]interface{}{
      "code":    400,
      "message": "invalid argument",
      "status":  "INVALID_ARGUMENT",
    },
  }

  testCases := []struct {
    name             string
    buildStubs       func(repo *mock.MockUserRepository)
    pathArgs         []interface{}
    body             interface{}
    expectedStatus   int
    expectedResponse interface{}
  }{
    {
      name: "OK",
      buildStubs: func(repo *mock.MockUserRepository) {
        repo.EXPECT().
          SetUserLocation(gomock.Any(), gomock.Eq(port.UserRepositorySetUserLocationRequest{
            Username: username,
            Point:    geo.Trunc(geo.Point{longitude, latitude}),
          })).
          Times(1).
          Return(port.UserRepositorySetUserLocationResponse{
            User: domain.User{
              Username: username,
            },
            PrevLocation: domain.Location{},
            Location: domain.Location{
              Point: geo.Point{longitude, latitude},
            },
          }, nil)
      },
      pathArgs: []interface{}{username},
      body: map[string]interface{}{
        "latitude":  latitude,
        "longitude": longitude,
      },
      expectedStatus: 200,
      expectedResponse: map[string]interface{}{
        "longitude": longitude,
        "latitude":  latitude,
      },
    },
    {
      name:             "no request body",
      buildStubs:       buildStubsNoCallExpected,
      pathArgs:         []interface{}{username},
      body:             nil,
      expectedStatus:   http.StatusBadRequest,
      expectedResponse: invalidArguentResponse,
    },
    {
      name:             "invalid request body",
      buildStubs:       buildStubsNoCallExpected,
      pathArgs:         []interface{}{username},
      body:             "invalid",
      expectedStatus:   http.StatusBadRequest,
      expectedResponse: invalidArguentResponse,
    },
    {
      name:       "empty username",
      buildStubs: buildStubsNoCallExpected,
      pathArgs:   []interface{}{""},
      body: map[string]interface{}{
        "latitude":  latitude,
        "longitude": longitude,
      },
      expectedStatus:   http.StatusBadRequest,
      expectedResponse: invalidArguentResponse,
    },
    {
      name:       "too shot username",
      buildStubs: buildStubsNoCallExpected,
      pathArgs:   []interface{}{shortUsername},
      body: map[string]interface{}{
        "latitude":  latitude,
        "longitude": longitude,
      },
      expectedStatus:   http.StatusBadRequest,
      expectedResponse: invalidArguentResponse,
    },
    {
      name:       "too long username",
      buildStubs: buildStubsNoCallExpected,
      pathArgs:   []interface{}{longUsername},
      body: map[string]interface{}{
        "latitude":  latitude,
        "longitude": longitude,
      },
      expectedStatus:   http.StatusBadRequest,
      expectedResponse: invalidArguentResponse,
    },
    {
      name:       "username contains a not allowed character",
      buildStubs: buildStubsNoCallExpected,
      pathArgs:   []interface{}{invalidUsername},
      body: map[string]interface{}{
        "latitude":  latitude,
        "longitude": longitude,
      },
      expectedStatus:   http.StatusBadRequest,
      expectedResponse: invalidArguentResponse,
    },
    {
      name:       "longitude less than min",
      buildStubs: buildStubsNoCallExpected,
      pathArgs:   []interface{}{username},
      body: map[string]interface{}{
        "latitude":  latitude,
        "longitude": -180.1,
      },
      expectedStatus:   http.StatusBadRequest,
      expectedResponse: invalidArguentResponse,
    },
    {
      name:       "longitude greater than max",
      buildStubs: buildStubsNoCallExpected,
      pathArgs:   []interface{}{username},
      body: map[string]interface{}{
        "latitude":  latitude,
        "longitude": 180.1,
      },
      expectedStatus:   http.StatusBadRequest,
      expectedResponse: invalidArguentResponse,
    },
    {
      name:       "latitude less than min",
      buildStubs: buildStubsNoCallExpected,
      pathArgs:   []interface{}{username},
      body: map[string]interface{}{
        "latitude":  -90.1,
        "longitude": longitude,
      },
      expectedStatus:   http.StatusBadRequest,
      expectedResponse: invalidArguentResponse,
    },
    {
      name:       "latitude greater than max",
      buildStubs: buildStubsNoCallExpected,
      pathArgs:   []interface{}{username},
      body: map[string]interface{}{
        "latitude":  90.1,
        "longitude": longitude,
      },
      expectedStatus:   http.StatusBadRequest,
      expectedResponse: invalidArguentResponse,
    },
    {
      name: "internal error",
      buildStubs: func(repo *mock.MockUserRepository) {
        repo.EXPECT().
          SetUserLocation(gomock.Any(), gomock.Any()).
          Times(1).
          Return(port.UserRepositorySetUserLocationResponse{}, fmt.Errorf("%w: %v", errpack.ErrInternalError, "test error"))
      },
      pathArgs: []interface{}{username},
      body: map[string]interface{}{
        "latitude":  latitude,
        "longitude": longitude,
      },
      expectedStatus: http.StatusInternalServerError,
      expectedResponse: map[string]interface{}{
        "error": map[string]interface{}{
          "code":    500,
          "message": "internal error",
          "status":  "INTERNAL",
        },
      },
    },
  }

  for _, tc := range testCases {
    tc := tc
    s.Run(tc.name, func() {
      ctrl := gomock.NewController(s.T())
      defer ctrl.Finish()

      logger := log.NewTestingLogger()

      repo := mock.NewMockUserRepository(ctrl)
      tc.buildStubs(repo)

      hc := mock.NewMockHistoryClient(ctrl)
      hc.EXPECT().AddRecord(gomock.Any(), gomock.Any()).AnyTimes()

      svc := service.NewUserService(repo, hc, logger)

      h := handler.NewHTTPHandler(svc, logger)

      server := httptest.NewServer(h)
      defer server.Close()

      e := httpexpect.New(s.T(), server.URL)

      req := e.PUT(path, tc.pathArgs...)
      //if tc.body != nil {
      req = req.WithJSON(tc.body).WithHeaders(
        map[string]string{
          "Content-Type": "application/json",
        },
      )
      //}

      res := req.Expect()

      res.Status(tc.expectedStatus)
      if tc.expectedResponse != nil {
        var b bytes.Buffer
        err := json.NewEncoder(&b).Encode(tc.expectedResponse)
        if err != nil {
          s.T().Fatal(err)
        }
        res.Body().Equal(b.String())
      } else {
        res.NoContent()
      }
    })
  }
}

func TestHTTPHandlerTestSuite(t *testing.T) {
  suite.Run(t, new(HTTPHandleTestSuite))
}
