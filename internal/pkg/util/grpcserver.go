package util

import (
	"context"
	"google.golang.org/grpc"
	"net"
	"time"
)

const (
	grpcStopTimeout = 30 * time.Second
)

// GRPCServer is a wrapper for grpc server.
// It provides Start and Stop methods.
type GRPCServer struct {
	server   *grpc.Server
	bindAddr string
}

// NewGRPCServer allocates and returns a new GRPCServer.
func NewGRPCServer(
	bindAddr string,
	registerServer func(*grpc.Server),
) *GRPCServer {
	server := grpc.NewServer()
	registerServer(server)

	return &GRPCServer{
		server:   server,
		bindAddr: bindAddr,
	}
}

// Start starts grpc server.
func (s *GRPCServer) Start() error {
	lis, err := net.Listen("tcp", s.bindAddr)
	if err != nil {
		return err
	}

	err = s.server.Serve(lis)
	if err != nil {
		return err
	}

	return nil
}

// Stop stops grpc server.
func (s *GRPCServer) Stop(ctx context.Context) error {
	stopped := make(chan struct{})
	go func() {
		s.server.GracefulStop()
		close(stopped)
	}()

	t := time.NewTimer(grpcStopTimeout)
	select {
	case <-t.C:
		s.server.Stop()
	case <-stopped:
		t.Stop()
	}

	return nil
}
