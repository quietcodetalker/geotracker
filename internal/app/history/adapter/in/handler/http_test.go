//go:build integration
// +build integration

package handler_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gavv/httpexpect"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gitlab.com/spacewalker/locations/internal/app/history/adapter/in/handler"
	"gitlab.com/spacewalker/locations/internal/app/history/core/port"
	"gitlab.com/spacewalker/locations/internal/app/history/core/port/mock"
	"gitlab.com/spacewalker/locations/internal/pkg/errpack"
	mocklog "gitlab.com/spacewalker/locations/internal/pkg/log/mock"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
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
	validFromStr := "2021-09-02T11:26:18+00:00"
	validFrom, err := time.Parse(time.RFC3339, validFromStr)
	require.NoError(s.T(), err)
	validToStr := "2022-09-02T11:26:18+00:00"
	validTo, err := time.Parse(time.RFC3339, validToStr)
	require.NoError(s.T(), err)
	username := "johnsmith"

	testCases := []struct {
		name             string
		urlParams        []interface{}
		queryParams      map[string]interface{}
		buildStubs       func(svc *mock.MockHistoryService)
		expectedStatus   int
		expectedResponse string
	}{
		{
			name:      "it returns INVALID_ARGUMENT if `from` query param contains invalid timestamp value",
			urlParams: []interface{}{"any"},
			queryParams: map[string]interface{}{
				"from": "invalid",
				"to":   "2021-09-02T11:26:18+00:00",
			},
			buildStubs: func(svc *mock.MockHistoryService) {
				svc.EXPECT().GetDistanceByUsername(gomock.Any(), gomock.Any()).Times(0)
			},
			expectedStatus: http.StatusBadRequest,
			expectedResponse: `
{
	"error": {
		"code": 400,
		"message": "invalid argument",
		"status": "INVALID_ARGUMENT"
	}
}
`,
		},
		{
			name:      "it returns INVALID_ARGUMENT if `to` query param contains invalid timestamp value",
			urlParams: []interface{}{"any"},
			queryParams: map[string]interface{}{
				"from": "2021-09-02T11:26:18+00:00",
				"to":   "invalid",
			},
			buildStubs: func(svc *mock.MockHistoryService) {
				svc.EXPECT().GetDistanceByUsername(gomock.Any(), gomock.Any()).Times(0)
			},
			expectedStatus: http.StatusBadRequest,
			expectedResponse: `
{
	"error": {
		"code": 400,
		"message": "invalid argument",
		"status": "INVALID_ARGUMENT"
	}
}
`,
		},
		{
			name:      "it passes both `to` and `from` query params to the service",
			urlParams: []interface{}{username},
			queryParams: map[string]interface{}{
				"from": validFromStr,
				"to":   validToStr,
			},
			buildStubs: func(svc *mock.MockHistoryService) {
				svc.EXPECT().
					GetDistanceByUsername(gomock.Any(), gomock.Eq(
						port.HistoryServiceGetDistanceByUsernameRequest{
							Username: username,
							From:     &validFrom,
							To:       &validTo,
						},
					)).
					Times(1)
			},
			expectedStatus:   0,
			expectedResponse: "",
		},
		{
			name:      "it passes both `to` query params  to the service if `from` missing",
			urlParams: []interface{}{username},
			queryParams: map[string]interface{}{
				"from": validFromStr,
			},
			buildStubs: func(svc *mock.MockHistoryService) {
				svc.EXPECT().
					GetDistanceByUsername(gomock.Any(), gomock.Eq(
						port.HistoryServiceGetDistanceByUsernameRequest{
							Username: username,
							From:     &validFrom,
							To:       nil,
						},
					)).
					Times(1)
			},
			expectedStatus:   0,
			expectedResponse: "",
		},
		{
			name:      "it passes both `from` query params  to the service if `to` missing",
			urlParams: []interface{}{username},
			queryParams: map[string]interface{}{
				"to": validToStr,
			},
			buildStubs: func(svc *mock.MockHistoryService) {
				svc.EXPECT().
					GetDistanceByUsername(gomock.Any(), gomock.Eq(
						port.HistoryServiceGetDistanceByUsernameRequest{
							Username: username,
							From:     nil,
							To:       &validTo,
						},
					)).
					Times(1)
			},
			expectedStatus:   0,
			expectedResponse: "",
		},
		{
			name:        "it does not pass both `from` and `to` to the service if both of them missing",
			urlParams:   []interface{}{username},
			queryParams: nil,
			buildStubs: func(svc *mock.MockHistoryService) {
				svc.EXPECT().
					GetDistanceByUsername(gomock.Any(), gomock.Eq(
						port.HistoryServiceGetDistanceByUsernameRequest{
							Username: username,
							From:     nil,
							To:       nil,
						},
					)).
					Times(1)
			},
			expectedStatus:   0,
			expectedResponse: "",
		},
		{
			name:        "it responds with `` if the service return `ErrInternalError`",
			urlParams:   []interface{}{username},
			queryParams: nil,
			buildStubs: func(svc *mock.MockHistoryService) {
				svc.EXPECT().
					GetDistanceByUsername(gomock.Any(), gomock.Any()).
					Times(1).
					Return(port.HistoryServiceGetDistanceByUsernameResponse{}, fmt.Errorf("%w", errpack.ErrInternalError))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedResponse: `
{
	"error": {
		"code": 500,
		"message": "internal error",
		"status": "INTERNAL"
	}
}
`,
		},
		{
			name:        "it responds with `` if the service return an unknown error",
			urlParams:   []interface{}{username},
			queryParams: nil,
			buildStubs: func(svc *mock.MockHistoryService) {
				svc.EXPECT().
					GetDistanceByUsername(gomock.Any(), gomock.Any()).
					Times(1).
					Return(port.HistoryServiceGetDistanceByUsernameResponse{}, errors.New("test error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedResponse: `
{
	"error": {
		"code": 500,
		"message": "unknown error",
		"status": "UNKNOWN"
	}
}
`,
		},
		{
			name:        "it responds with `FAILED_PRECONDITION` if the service return `ErrFailedPrecondition`",
			urlParams:   []interface{}{username},
			queryParams: nil,
			buildStubs: func(svc *mock.MockHistoryService) {
				svc.EXPECT().
					GetDistanceByUsername(gomock.Any(), gomock.Any()).
					Times(1).
					Return(
						port.HistoryServiceGetDistanceByUsernameResponse{},
						fmt.Errorf(
							"%w: %v",
							errpack.ErrFailedPrecondition,
							errors.New("test error"),
						),
					)
			},
			expectedStatus: http.StatusUnprocessableEntity,
			expectedResponse: `
{
	"error": {
		"code": 422,
		"message": "failed precondition: test error",
		"status": "FAILED_PRECONDITION"
	}
}
`,
		},
		{
			name:        "it responds with `INVALID_ARGUMENT` if the service return `ErrInvalidArgument`",
			urlParams:   []interface{}{username},
			queryParams: nil,
			buildStubs: func(svc *mock.MockHistoryService) {
				svc.EXPECT().
					GetDistanceByUsername(gomock.Any(), gomock.Any()).
					Times(1).
					Return(
						port.HistoryServiceGetDistanceByUsernameResponse{},
						fmt.Errorf(
							"%w: %v",
							errpack.ErrInvalidArgument,
							errors.New("test error"),
						),
					)
			},
			expectedStatus: http.StatusBadRequest,
			expectedResponse: `
{
	"error": {
		"code": 400,
		"message": "invalid argument: test error",
		"status": "INVALID_ARGUMENT"
	}
}
`,
		},
		{
			name:        "it responds with `NOT_FOUND` if the service return `ErrNotFound`",
			urlParams:   []interface{}{username},
			queryParams: nil,
			buildStubs: func(svc *mock.MockHistoryService) {
				svc.EXPECT().
					GetDistanceByUsername(gomock.Any(), gomock.Any()).
					Times(1).
					Return(
						port.HistoryServiceGetDistanceByUsernameResponse{},
						fmt.Errorf(
							"%w: %v",
							errpack.ErrNotFound,
							errors.New("test error"),
						),
					)
			},
			expectedStatus: http.StatusNotFound,
			expectedResponse: `
{
	"error": {
		"code": 404,
		"message": "not found: test error",
		"status": "NOT_FOUND"
	}
}
`,
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.T().Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(s.T())

			svc := mock.NewMockHistoryService(ctrl)
			logger := mocklog.NewMockLogger(ctrl)

			logger.EXPECT().Info(gomock.Any(), gomock.Any()) // Ignore logging

			tc.buildStubs(svc)

			h := handler.NewHTTPHandler(svc, logger)

			server := httptest.NewServer(h)
			defer server.Close()

			e := httpexpect.New(s.T(), server.URL)

			req := e.GET("/users/{username}/distance", tc.urlParams...)
			for k, v := range tc.queryParams {
				req = req.WithQuery(k, v)
			}

			response := req.Expect()
			if tc.expectedStatus != 0 {
				response = response.Status(tc.expectedStatus)
			}
			if tc.expectedResponse != "" {
				var obj interface{}
				err := json.Unmarshal([]byte(tc.expectedResponse), &obj)
				require.NoError(t, err)

				response.JSON().Equal(obj)
			}
		})
	}
}
