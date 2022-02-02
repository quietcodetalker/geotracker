package location

import (
	"context"
	"errors"
	"fmt"
	"gitlab.com/spacewalker/locations/internal/app/location/adapter/in/handler"
	"gitlab.com/spacewalker/locations/internal/app/location/adapter/out/historyclient"
	"gitlab.com/spacewalker/locations/internal/app/location/adapter/out/repository"
	"gitlab.com/spacewalker/locations/internal/app/location/core/service"
	"gitlab.com/spacewalker/locations/internal/pkg/config"
	"gitlab.com/spacewalker/locations/internal/pkg/log"
	"gitlab.com/spacewalker/locations/internal/pkg/middleware"
	"gitlab.com/spacewalker/locations/internal/pkg/util"
	pb "gitlab.com/spacewalker/locations/pkg/api/proto/v1/location"
	"google.golang.org/grpc"
	log2 "log"
	"strings"
	"sync"
)

// App is a history application.
type App struct {
	config     config.LocationConfig
	httpServer *util.HTTPServer
	grpcServer *util.GRPCServer
	logger     log.Logger
}

// NewApp creates and instance of location application and returns its pointer.
func NewApp(config config.LocationConfig) *App {
	return &App{
		config: config,
	}
}

// Start starts the application.
func (a *App) Start() error {
	dbSource := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		a.config.DBHost, a.config.DBPort, a.config.DBUser, a.config.DBPassword, a.config.DBName, a.config.DBSSLMode,
	)
	db, err := util.OpenDB(a.config.DBDriver, dbSource)
	if err != nil {
		return fmt.Errorf("failed to open db: %v", err)
	}

	a.logger, err = log.NewZapLogger(true)
	if err != nil {
		log2.Panic()
	}

	repo := repository.NewPostgresRepository(db)
	historyClient := historyclient.NewGRPCClient(a.config.HistoryAddr)
	svc := service.NewUserService(repo, historyClient, a.logger)
	httpHandler := handler.NewHTTPHandler(svc, a.logger)
	grpcHandler := handler.NewGRPCInternalHandler(svc)

	a.httpServer = util.NewHTTPServer(a.config.BindAddrHTTP, httpHandler)
	a.grpcServer = util.NewGRPCServer(
		a.config.BindAddrGRPC,
		func(server *grpc.Server) {
			pb.RegisterLocationInternalServer(server, grpcHandler)
		},
		grpc.UnaryInterceptor(middleware.LoggerUnaryServerInterceptor(a.logger)),
	)

	var httpErr, grpcErr error

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		a.logger.Info(fmt.Sprintf("Starting HTTP server on %v", a.config.BindAddrHTTP), nil)
		if httpErr = a.httpServer.Start(); httpErr != nil {
			a.logger.Info(fmt.Sprintf("failed to start http server: %v", httpErr), nil)
		}
		wg.Done()
	}()

	go func() {
		//log.Printf("Starting gRPC server on %v", a.config.BindAddrGRPC)
		a.logger.Info(fmt.Sprintf("Starting gRPC server on %v", a.config.BindAddrGRPC), nil)
		if grpcErr = a.grpcServer.Start(); grpcErr != nil {
			a.logger.Error(fmt.Sprintf("failed to start grpc server: %v", grpcErr), nil)
		}
		wg.Done()
	}()

	wg.Wait()

	errMsgs := make([]string, 0, 2)

	if grpcErr != nil {
		errMsgs = append(errMsgs, fmt.Sprintf("failed to start grpc server: %v", grpcErr))
	}
	if httpErr != nil {
		errMsgs = append(errMsgs, fmt.Sprintf("failed to start http server: %v", httpErr))
	}

	if len(errMsgs) > 0 {
		return errors.New(strings.Join(errMsgs, ""))
	}

	return nil
}

// Stop stops the application.
func (a *App) Stop(ctx context.Context) error {
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		if err := a.httpServer.Stop(context.Background()); err != nil {
			a.logger.Error(fmt.Sprintf("failed to stop HTTP server :%v", err), nil)
		}
		wg.Done()
	}()
	go func() {
		if err := a.grpcServer.Stop(context.Background()); err != nil {
			a.logger.Error(fmt.Sprintf("failed to gRPC http server :%v", err), nil)
		}
		wg.Done()
	}()

	wg.Wait()

	return nil
}
