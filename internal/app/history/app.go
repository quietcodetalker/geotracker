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
	"gitlab.com/spacewalker/locations/internal/pkg/util"
	pb "gitlab.com/spacewalker/locations/pkg/api/proto/v1/history"
	"google.golang.org/grpc"
	"log"
	"strings"
	"sync"
)

// App is a history application.
type App struct {
	config     config.HistoryConfig
	httpServer *util.HTTPServer
	grpcServer *util.GRPCServer
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

	repo := repository.NewPostgresRepository(db)
	locationClient := locationclient.NewGRPCClient(a.config.LocationAddr)
	svc := service.NewHistoryService(repo, locationClient)
	httpHandler := handler.NewHTTPHandler(svc)
	grpcHandler := handler.NewGRPCHandler(svc)

	a.httpServer = util.NewHTTPServer(a.config.BindAddrHTTP, httpHandler)
	a.grpcServer = util.NewGRPCServer(a.config.BindAddrGRPC, func(server *grpc.Server) {
		pb.RegisterHistoryServer(server, grpcHandler)
	})

	var httpErr, grpcErr error

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		log.Printf("Starting HTTP server on %v", a.config.BindAddrHTTP)
		if httpErr = a.httpServer.Start(); httpErr != nil {
			log.Printf("failed to start http server: %v", httpErr)
		}
		wg.Done()
	}()

	go func() {
		log.Printf("Starting gRPC server on %v", a.config.BindAddrGRPC)
		if grpcErr = a.grpcServer.Start(); grpcErr != nil {
			log.Printf("failed to start grpc server: %v", grpcErr)
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
			log.Printf("failed to stop http server :%v", err)
		}
		wg.Done()
	}()
	go func() {
		if err := a.grpcServer.Stop(context.Background()); err != nil {
			log.Printf("failed to stop grpc server :%v", err)
		}
		wg.Done()
	}()

	wg.Wait()

	return nil
}
