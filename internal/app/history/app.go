package history

import (
	"context"
	"errors"
	"fmt"
	"gitlab.com/spacewalker/locations/internal/app/history/adapter/in/handler"
	"gitlab.com/spacewalker/locations/internal/app/history/adapter/out/locationclient"
	"gitlab.com/spacewalker/locations/internal/app/history/adapter/out/repository"
	"gitlab.com/spacewalker/locations/internal/app/history/core/service"
	"gitlab.com/spacewalker/locations/internal/pkg/config"
	"gitlab.com/spacewalker/locations/internal/pkg/log"
	"gitlab.com/spacewalker/locations/internal/pkg/util"
	pb "gitlab.com/spacewalker/locations/pkg/api/proto/v1/history"
	"google.golang.org/grpc"
	log2 "log"
	"strings"
	"sync"
)

// App is a history application.
type App struct {
	config     config.HistoryConfig
	httpServer *util.HTTPServer
	grpcServer *util.GRPCServer
	logger     log.Logger
}

// NewApp creates and instance of history application and returns its pointer.
func NewApp(config config.HistoryConfig) *App {
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
		log2.Panic(err)
	}

	repo := repository.NewPostgresRepository(db)
	locationClient := locationclient.NewGRPCClient(a.config.LocationAddr)
	svc := service.NewHistoryService(repo, locationClient, a.logger)
	httpHandler := handler.NewHTTPHandler(svc, a.logger)
	grpcHandler := handler.NewGRPCHandler(svc)

	a.httpServer = util.NewHTTPServer(a.config.BindAddrHTTP, httpHandler)
	a.grpcServer = util.NewGRPCServer(a.config.BindAddrGRPC, func(server *grpc.Server) {
		pb.RegisterHistoryServer(server, grpcHandler)
	})

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
		a.logger.Info(fmt.Sprintf("Starting gRPC server on %v", a.config.BindAddrGRPC), nil)
		if grpcErr = a.grpcServer.Start(); grpcErr != nil {
			a.logger.Info(fmt.Sprintf("failed to start grpc server: %v", grpcErr), nil)
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
			a.logger.Info(fmt.Sprintf("failed to stop http server :%v", err), nil)
		}
		wg.Done()
	}()
	go func() {
		if err := a.grpcServer.Stop(context.Background()); err != nil {
			a.logger.Info(fmt.Sprintf("failed to stop grpc server :%v", err), nil)
		}
		wg.Done()
	}()

	wg.Wait()

	return nil
}
