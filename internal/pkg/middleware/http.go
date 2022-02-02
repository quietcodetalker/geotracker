package middleware

import (
	"gitlab.com/spacewalker/locations/internal/pkg/errpack"
	"gitlab.com/spacewalker/locations/internal/pkg/log"
	"gitlab.com/spacewalker/locations/internal/pkg/util"
	"net/http"
	"time"
)

type responseData struct {
	status int
	size   int
}

type loggingResponseWriter struct {
	http.ResponseWriter
	responseData *responseData
}

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

func LoggerMiddleware(logger log.Logger) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			responseData := &responseData{
				status: 0,
				size:   0,
			}

			lw := loggingResponseWriter{
				ResponseWriter: w,
				responseData:   responseData,
			}

			uri := r.RequestURI
			method := r.Method

			next.ServeHTTP(&lw, r)

			duration := time.Since(start)

			logger.Info("http request complete", log.Fields{
				"uri":      uri,
				"method":   method,
				"duration": duration,
				"status":   responseData.status,
				"size":     responseData.size,
			})
		})
	}
}

func RecovererMiddleware(logger log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			uri := r.RequestURI
			method := r.Method

			defer func() {
				err := recover()
				if err != nil {
					logger.Error("panic recovered", log.Fields{
						"uri":    uri,
						"method": method,
						"error":  err,
					})

					status, body := errpack.ErrToHTTP(errpack.ErrInternalError)
					util.Respond(w, status, body)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
