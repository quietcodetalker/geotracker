package history

import (
	"fmt"
	"gitlab.com/spacewalker/locations/internal/app/history/adapter/in/handler"
	"gitlab.com/spacewalker/locations/internal/app/history/adapter/out/repository"
	"gitlab.com/spacewalker/locations/internal/app/history/core/service"
	"gitlab.com/spacewalker/locations/internal/pkg/config"
	"gitlab.com/spacewalker/locations/internal/pkg/util"
	pb "gitlab.com/spacewalker/locations/pkg/api/proto/v1/history"
	"google.golang.org/grpc"
	"log"
	"net"
)

// App is a history application.
type App struct {
	config config.HistoryConfig
	server *grpc.Server
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
	svc := service.NewHistoryService(repo)
	hdl := handler.NewGRPCHandler(svc)

	log.Println("Listening on " + a.config.BindAddr)
	lis, err := net.Listen("tcp", a.config.BindAddr)
	if err != nil {
		return err
	}

	a.server = grpc.NewServer()

	pb.RegisterHistoryServer(a.server, hdl)

	if err := a.server.Serve(lis); err != nil {
		if err == grpc.ErrServerStopped {
			return nil
		}

		return err
	}

	return nil
}

// Stop stops the application.
func (a *App) Stop() error {
	a.server.GracefulStop()

	return nil
}
