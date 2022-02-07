package handler_test

import (
	"fmt"
	"time"

	"github.com/golang/mock/gomock"
	"gitlab.com/spacewalker/locations/internal/app/history/core/port"
)

type eqHistoryServiceGetDistanceByUsernameRequestMatcher struct {
	req port.HistoryServiceGetDistanceByUsernameRequest
}

func (m eqHistoryServiceGetDistanceByUsernameRequestMatcher) Matches(x interface{}) bool {
	req, ok := x.(port.HistoryServiceGetDistanceByUsernameRequest)
	if !ok {
		return false
	}

	if m.req.Username != req.Username {
		return false
	}

	if m.req.From == nil && req.From != nil ||
		m.req.From != nil && req.From == nil ||
		m.req.To == nil && req.To != nil ||
		m.req.To != nil && req.To == nil {
		return false
	}

	if m.req.From != nil {
		diff := m.req.From.Sub(*req.From)
		if diff < 0 {
			diff = -diff
		}
		if diff > time.Second {
			return false
		}
	}

	if m.req.To != nil {
		diff := m.req.To.Sub(*req.To)
		if diff < 0 {
			diff = -diff
		}
		if diff > time.Second {
			return false
		}
	}

	return true
}

func EqHistoryServiceGetDistanceByUsernameRequest(req port.HistoryServiceGetDistanceByUsernameRequest) gomock.Matcher {
	return eqHistoryServiceGetDistanceByUsernameRequestMatcher{
		req: req,
	}
}

func (m eqHistoryServiceGetDistanceByUsernameRequestMatcher) String() string {
	return fmt.Sprintf("matches req %v", m.req)
}
