package location

import (
	"context"
	"errors"
	"fmt"
	log2 "log"
	"strings"
	"sync"
	"time"

	"gitlab.com/spacewalker/geotracker/internal/pkg/retrier"

	"github.com/go-chi/chi/v5"
	"github.com/sony/gobreaker"
	"gitlab.com/spacewalker/geotracker/internal/app/location/adapter/in/handler"
	"gitlab.com/spacewalker/geotracker/internal/app/location/adapter/out/historyclient"
	"gitlab.com/spacewalker/geotracker/internal/app/location/adapter/out/repository"
	"gitlab.com/spacewalker/geotracker/internal/app/location/core/service"
	"gitlab.com/spacewalker/geotracker/internal/pkg/config"
	"gitlab.com/spacewalker/geotracker/internal/pkg/errpack"
	"gitlab.com/spacewalker/geotracker/internal/pkg/log"
	"gitlab.com/spacewalker/geotracker/internal/pkg/middleware"
	"gitlab.com/spacewalker/geotracker/internal/pkg/util"
	pb "gitlab.com/spacewalker/geotracker/pkg/api/proto/v1/location"
	"google.golang.org/grpc"
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

	a.logger, err = log.NewZapLogger(a.config.AppEnv == "development")
	if err != nil {
		log2.Panic()
	}

	cb := gobreaker.NewCircuitBreaker(gobreaker.Settings{
		Name:        "historyclient",
		MaxRequests: 3,
		Interval:    5 * time.Second,
		Timeout:     7 * time.Second,
		IsSuccessful: func(err error) bool {
			return !errors.Is(err, errpack.ErrInternalError)
		},
	})
	re := retrier.New(retrier.Config{
		Delay:   3 * time.Second,
		Retries: 3,
		IsSuccessful: func(err error) bool {
			return !errors.Is(err, errpack.ErrInternalError)
		},
	})

	repo := repository.NewPostgresRepository(db)
	historyClient := historyclient.NewGRPCClient(a.config.HistoryAddr, a.logger)
	proxifiedHistoryClient := historyclient.NewProxy(historyClient, cb, re)
	svc := service.NewUserService(repo, proxifiedHistoryClient, a.logger)
	httpHandler := handler.NewHTTPHandler(svc, a.logger)
	grpcHandler := handler.NewGRPCHandler(svc)

	rootHandler := chi.NewRouter()
	rootHandler.Mount("/v1", httpHandler)

	a.httpServer = util.NewHTTPServer(a.config.BindAddrHTTP, rootHandler)
	a.grpcServer = util.NewGRPCServer(
		a.config.BindAddrGRPC,
		func(server *grpc.Server) {
			pb.RegisterLocationInternalServer(server, grpcHandler)
		},
		grpc.ChainUnaryInterceptor(
			middleware.TracingUnaryServerInterceptor(a.logger),
			middleware.LoggerUnaryServerInterceptor(a.logger),
		),
	)

	var httpErr, grpcErr error

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		a.logger.Info(fmt.Sprintf("Starting HTTP server on %v", a.config.BindAddrHTTP), nil)
		if httpErr = a.httpServer.Start(); httpErr != nil {
			a.logger.Error(fmt.Sprintf("failed to start http server: %v", httpErr), nil)
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

	cancelCtx, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()

	go func() {
		if err := a.httpServer.Stop(cancelCtx); err != nil {
			a.logger.Error(fmt.Sprintf("failed to stop HTTP server :%v", err), nil)
		}
		wg.Done()
	}()
	go func() {
		if err := a.grpcServer.Stop(cancelCtx); err != nil {
			a.logger.Error(fmt.Sprintf("failed to gRPC http server :%v", err), nil)
		}
		wg.Done()
	}()

	wg.Wait()

	return nil
}
