package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"gitlab.com/spacewalker/geotracker/internal/pkg/errpack"

	"github.com/gavv/httpexpect"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gitlab.com/spacewalker/geotracker/internal/app/history/adapter/in/handler"
	"gitlab.com/spacewalker/geotracker/internal/app/history/core/port"
	"gitlab.com/spacewalker/geotracker/internal/app/history/core/port/mock"
	mocklog "gitlab.com/spacewalker/geotracker/internal/pkg/log/mock"
	"gitlab.com/spacewalker/geotracker/internal/pkg/util"
)

// HistoryHTTPHandlerTestSuite is a test suite that covers history http handler functionality.
type HistoryHTTPHandlerTestSuite struct {
	suite.Suite
}

func (s *HistoryHTTPHandlerTestSuite) SetupTest() {
}

func (s *HistoryHTTPHandlerTestSuite) TearDownTest() {
}

func (s *HistoryHTTPHandlerTestSuite) TearDownSuite() {
}

func TestUserHTTPHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(HistoryHTTPHandlerTestSuite))
}

func (s *HistoryHTTPHandlerTestSuite) Test_GetDistance() {
	invalidUsernameCharacterSet := "!@#$%^&*()_+=-'\""
	getUserDistancePath := "/users/{validUsername}/distance"
	validUsername := util.RandomUsername()
	util.RandomString(1, 1, invalidUsernameCharacterSet)
	from, to := util.RandomTimeInterval()
	validFromStr := from.Format("2006-01-02T15:04:05-07:00")
	validToStr := to.Format("2006-01-02T15:04:05-07:00")
	distance := util.RandomFloat64(0.0, 1000.0)
	invaildTimestampStr := "invalid"
	customErrMsg := util.RandomString(4, 10, util.CharacterSetAlphanumeric)

	testCases := []struct {
		name             string
		path             string
		urlParams        []interface{}
		queryParams      map[string]interface{}
		headers          map[string]string
		buildStubs       func(service *mock.MockHistoryService)
		expectedStatus   int
		expectedResponse interface{}
	}{
		{
			name:      "it responds with OK if valid `validUsername`, `from` and `to` are provided",
			path:      getUserDistancePath,
			urlParams: []interface{}{validUsername},
			queryParams: map[string]interface{}{
				"from": validFromStr,
				"to":   validToStr,
			},
			headers: map[string]string{
				"Content-Type": "application/json",
			},
			buildStubs: func(svc *mock.MockHistoryService) {
				svc.EXPECT().
					GetDistanceByUsername(
						gomock.Any(),
						EqHistoryServiceGetDistanceByUsernameRequest(port.HistoryServiceGetDistanceByUsernameRequest{
							Username: validUsername,
							From:     &from,
							To:       &to,
						}),
					).
					Times(1).
					Return(port.HistoryServiceGetDistanceByUsernameResponse{
						Distance: distance,
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedResponse: port.HistoryServiceGetDistanceByUsernameResponse{
				Distance: distance,
			},
		},
		{
			name:      "it responds with OK if valid `validUsername` and `from` are provided without `to`",
			path:      getUserDistancePath,
			urlParams: []interface{}{validUsername},
			queryParams: map[string]interface{}{
				"from": validFromStr,
			},
			headers: map[string]string{
				"Content-Type": "application/json",
			},
			buildStubs: func(svc *mock.MockHistoryService) {
				svc.EXPECT().
					GetDistanceByUsername(
						gomock.Any(),
						EqHistoryServiceGetDistanceByUsernameRequest(port.HistoryServiceGetDistanceByUsernameRequest{
							Username: validUsername,
							From:     &from,
							To:       nil,
						}),
					).
					Times(1).
					Return(port.HistoryServiceGetDistanceByUsernameResponse{
						Distance: distance,
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedResponse: port.HistoryServiceGetDistanceByUsernameResponse{
				Distance: distance,
			},
		},
		{
			name:      "it responds with OK if valid `validUsername` and `to` are provided without `from`",
			path:      getUserDistancePath,
			urlParams: []interface{}{validUsername},
			queryParams: map[string]interface{}{
				"to": validToStr,
			},
			headers: map[string]string{
				"Content-Type": "application/json",
			},
			buildStubs: func(svc *mock.MockHistoryService) {
				svc.EXPECT().
					GetDistanceByUsername(
						gomock.Any(),
						EqHistoryServiceGetDistanceByUsernameRequest(port.HistoryServiceGetDistanceByUsernameRequest{
							Username: validUsername,
							From:     nil,
							To:       &to,
						}),
					).
					Times(1).
					Return(port.HistoryServiceGetDistanceByUsernameResponse{
						Distance: distance,
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedResponse: port.HistoryServiceGetDistanceByUsernameResponse{
				Distance: distance,
			},
		},
		{
			name:        "it responds with OK if valid `validUsername` is provided without `to` and `from`",
			path:        getUserDistancePath,
			urlParams:   []interface{}{validUsername},
			queryParams: nil,
			headers: map[string]string{
				"Content-Type": "application/json",
			},
			buildStubs: func(svc *mock.MockHistoryService) {
				svc.EXPECT().
					GetDistanceByUsername(
						gomock.Any(),
						EqHistoryServiceGetDistanceByUsernameRequest(port.HistoryServiceGetDistanceByUsernameRequest{
							Username: validUsername,
							From:     nil,
							To:       nil,
						}),
					).
					Times(1).
					Return(port.HistoryServiceGetDistanceByUsernameResponse{
						Distance: distance,
					}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedResponse: port.HistoryServiceGetDistanceByUsernameResponse{
				Distance: distance,
			},
		},
		{
			name:      "it responds with BAD_REQUEST if invalid `from` is provided",
			path:      getUserDistancePath,
			urlParams: []interface{}{validUsername},
			queryParams: map[string]interface{}{
				"from": invaildTimestampStr,
			},
			headers: map[string]string{
				"Content-Type": "application/json",
			},
			buildStubs: func(svc *mock.MockHistoryService) {
				svc.EXPECT().
					GetDistanceByUsername(
						gomock.Any(),
						gomock.Any(),
					).
					Times(0)
			},
			expectedStatus: http.StatusBadRequest,
			expectedResponse: map[string]interface{}{
				"error": map[string]interface{}{
					"code":    400,
					"message": "invalid argument",
					"status":  "INVALID_ARGUMENT",
				},
			},
		},
		{
			name:      "it responds with BAD_REQUEST if invalid `to` is provided",
			path:      getUserDistancePath,
			urlParams: []interface{}{validUsername},
			queryParams: map[string]interface{}{
				"to": invaildTimestampStr,
			},
			headers: map[string]string{
				"Content-Type": "application/json",
			},
			buildStubs: func(svc *mock.MockHistoryService) {
				svc.EXPECT().
					GetDistanceByUsername(
						gomock.Any(),
						gomock.Any(),
					).
					Times(0)
			},
			expectedStatus: http.StatusBadRequest,
			expectedResponse: map[string]interface{}{
				"error": map[string]interface{}{
					"code":    400,
					"message": "invalid argument",
					"status":  "INVALID_ARGUMENT",
				},
			},
		},
		{
			name:        "it responds with NOT_FOUND if service returns ErrNotFound",
			path:        getUserDistancePath,
			urlParams:   []interface{}{validUsername},
			queryParams: nil,
			headers: map[string]string{
				"Content-Type": "application/json",
			},
			buildStubs: func(svc *mock.MockHistoryService) {
				svc.EXPECT().
					GetDistanceByUsername(
						gomock.Any(),
						EqHistoryServiceGetDistanceByUsernameRequest(port.HistoryServiceGetDistanceByUsernameRequest{
							Username: validUsername,
							From:     nil,
							To:       nil,
						}),
					).
					Times(1).
					Return(
						port.HistoryServiceGetDistanceByUsernameResponse{
							Distance: distance,
						},
						fmt.Errorf("%w: %v", errpack.ErrNotFound, errors.New(customErrMsg)),
					)
			},
			expectedStatus: http.StatusNotFound,
			expectedResponse: map[string]interface{}{
				"error": map[string]interface{}{
					"code":    404,
					"message": fmt.Sprintf("%s: %s", errpack.ErrNotFound.Error(), customErrMsg),
					"status":  "NOT_FOUND",
				},
			},
		},
		{
			name:        "it responds with INTERNAL if service returns ErrInternalError",
			path:        getUserDistancePath,
			urlParams:   []interface{}{validUsername},
			queryParams: nil,
			headers: map[string]string{
				"Content-Type": "application/json",
			},
			buildStubs: func(svc *mock.MockHistoryService) {
				svc.EXPECT().
					GetDistanceByUsername(
						gomock.Any(),
						EqHistoryServiceGetDistanceByUsernameRequest(port.HistoryServiceGetDistanceByUsernameRequest{
							Username: validUsername,
							From:     nil,
							To:       nil,
						}),
					).
					Times(1).
					Return(
						port.HistoryServiceGetDistanceByUsernameResponse{
							Distance: distance,
						},
						fmt.Errorf("%w: %v", errpack.ErrInternalError, errors.New(customErrMsg)),
					)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedResponse: map[string]interface{}{
				"error": map[string]interface{}{
					"code":    500,
					"message": errpack.ErrInternalError.Error(),
					"status":  "INTERNAL",
				},
			},
		},
		{
			name:        "it responds with UNKNOWN if service returns an unknown error",
			path:        getUserDistancePath,
			urlParams:   []interface{}{validUsername},
			queryParams: nil,
			headers: map[string]string{
				"Content-Type": "application/json",
			},
			buildStubs: func(svc *mock.MockHistoryService) {
				svc.EXPECT().
					GetDistanceByUsername(
						gomock.Any(),
						EqHistoryServiceGetDistanceByUsernameRequest(port.HistoryServiceGetDistanceByUsernameRequest{
							Username: validUsername,
							From:     nil,
							To:       nil,
						}),
					).
					Times(1).
					Return(
						port.HistoryServiceGetDistanceByUsernameResponse{
							Distance: distance,
						},
						fmt.Errorf("%v", errors.New(customErrMsg)),
					)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedResponse: map[string]interface{}{
				"error": map[string]interface{}{
					"code":    500,
					"message": "unknown error",
					"status":  "UNKNOWN",
				},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.Run(tc.name, func() {
			ctrl := gomock.NewController(s.T())

			svc := mock.NewMockHistoryService(ctrl)
			logger := mocklog.NewMockLogger(ctrl)

			logger.EXPECT().Info(gomock.Any(), gomock.Any()) // Ignore logging

			tc.buildStubs(svc)

			h := handler.NewHTTPHandler(svc, logger)

			server := httptest.NewServer(h)
			defer server.Close()

			e := httpexpect.New(s.T(), server.URL)

			req := e.GET(tc.path, tc.urlParams...).WithHeaders(tc.headers)

			for k, v := range tc.queryParams {
				req = req.WithQuery(k, v)
			}

			res := req.Expect()

			res.Header("Content-Type").Equal("application/json")

			res.Status(tc.expectedStatus)
			if tc.expectedResponse != nil {
				var b bytes.Buffer
				err := json.NewEncoder(&b).Encode(tc.expectedResponse)
				require.NoError(s.T(), err)
				res.Body().Equal(b.String())
			} else {
				res.NoContent()
			}
		})
	}
}
