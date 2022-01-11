package main

import (
	"fmt"
	"gitlab.com/spacewalker/locations/internal/app/location/adapter/in/handler"
	"gitlab.com/spacewalker/locations/internal/app/location/adapter/out/repository"
	"gitlab.com/spacewalker/locations/internal/app/location/core/service"
	"gitlab.com/spacewalker/locations/internal/pkg/config"
	"gitlab.com/spacewalker/locations/internal/pkg/util"
	"log"
	"path"
	"runtime"
)

func main() {
	_, filename, _, _ := runtime.Caller(0)
	rootDir := path.Join(path.Dir(filename), "../..")

	cfg, err := config.LoadUserConfig(
		"user",
		path.Join(rootDir, "configs"),
	)
	if err != nil {
		log.Panicf("failed to load config: %v", err)
	}

	dbSource := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode,
	)
	db, err := util.OpenDB(cfg.DBDriver, dbSource)
	if err != nil {
		log.Panicf("failed to open db: %v", err)
	}

	repo := repository.NewPostgresRepository(db)
	svc := service.NewUserService(repo)

	grpcHandler := handler.NewGRPCHandler(svc)
	if err := grpcHandler.Start(":50051"); err != nil {
		log.Panic(err)
	}
}
