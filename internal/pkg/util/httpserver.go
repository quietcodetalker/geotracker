package util

import (
	"context"
	"log"
	"net/http"
	"time"
)

const (
	httpStopTimeout = 30 * time.Second
)

// HTTPServer is a http server wrapper that provides Stop and Start methods.
type HTTPServer struct {
	server *http.Server
}

// NewHTTPServer allocates and returns a new HTTPServer.
func NewHTTPServer(bindAddr string, handler http.Handler) *HTTPServer {
	return &HTTPServer{
		server: &http.Server{
			Addr:    bindAddr,
			Handler: handler,
		},
	}
}

// Start starts http server.
func (s *HTTPServer) Start() error {
	if err := s.server.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			return nil
		}
		return err
	}

	return nil
}

// Stop stops http server
func (s *HTTPServer) Stop(ctx context.Context) error {
	stopCtx, cancelStopCtx := context.WithTimeout(ctx, httpStopTimeout)
	defer cancelStopCtx()

	go func() {
		<-stopCtx.Done()
		if stopCtx.Err() == context.DeadlineExceeded {
			log.Println("graceful shutdown of the server is timed out...forcing exit.")
		}
	}()

	if err := s.server.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}
