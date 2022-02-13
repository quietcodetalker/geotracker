package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"path"
	"runtime"
	"syscall"

	"gitlab.com/spacewalker/geotracker/internal/app/location"
	"gitlab.com/spacewalker/geotracker/internal/pkg/config"
)

func main() {
	_, filename, _, _ := runtime.Caller(0)
	rootDir := path.Join(path.Dir(filename), "../..")

	cfg, err := config.LoadLocationConfig(
		"locations",
		path.Join(rootDir, "configs"),
	)
	if err != nil {
		log.Panicf("failed to load config: %v", err)
	}

	app := location.NewApp(cfg)

	if err := app.Start(); err != nil {
		log.Panic(err)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	<-sig
	app.Stop(context.Background())
}
