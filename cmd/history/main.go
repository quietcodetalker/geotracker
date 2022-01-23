package main

import (
	"gitlab.com/spacewalker/locations/internal/app/history"
	"gitlab.com/spacewalker/locations/internal/pkg/config"
	"log"
	"path"
	"runtime"
)

func main() {
	_, filename, _, _ := runtime.Caller(0)
	rootDir := path.Join(path.Dir(filename), "../..")

	cfg, err := config.LoadHistoryConfig(
		"history",
		path.Join(rootDir, "configs"),
	)
	if err != nil {
		log.Panicf("failed to load config: %v", err)
	}

	app := history.NewApp(cfg)

	if err := app.Start(); err != nil {
		log.Panic(err)
	}
}
