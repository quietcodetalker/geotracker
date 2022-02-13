package util

import (
	"context"
	"errors"

	"gitlab.com/spacewalker/geotracker/internal/pkg/errpack"
	"gitlab.com/spacewalker/geotracker/internal/pkg/log"
)

// LogInternalError logs error with loggin level ERROR in case of ErrInternalError
func LogInternalError(ctx context.Context, logger log.Logger, err error, args ...interface{}) {
	if errors.Is(err, errpack.ErrInternalError) {
		traceID, _ := GetTraceIDFromCtx(ctx)
		logger.Error(err.Error(), log.Fields{
			"trace-id": traceID,
			"args":     args,
		})
	}
}
