package main

import (
	"context"
	_ "github.com/lib/pq"
	"neatly"
	"neatly/internal/handler"
	"neatly/internal/repository"
	"neatly/internal/service"
	"neatly/internal/session"
	"neatly/pkg/logging"
	"neatly/pkg/postgres"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	logging.Init()
	logger := logging.GetLogger()

	cfg := session.GetConfig()

	db, err := postgres.NewDB(cfg.DB)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Info("initializing repository")
	repos := repository.New(db, logger)

	logger.Info("initializing services")
	services := service.New(repos, logger)

	logger.Info("initializing handler")
	handlers := handler.New(services, logger)

	srv := new(neatly.Server)
	go func() {
		if err := srv.Run(cfg.Port, handlers.RegisterHandler(cfg.IsDebug)); err != nil {
			logger.Fatal(err)
		}
	}()
	logger.Infof("Server is running on %v", cfg.Port)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGTERM, syscall.SIGINT)
	<-shutdown

	logger.Info("Shutting down...")
	if err := srv.ShutdownGraceful(context.Background()); err != nil {
		logger.Errorf("error occured while shutting down: %s", err.Error())
	}

	if err := db.Close(); err != nil {
		logger.Errorf("error occured while closing database connection: %s", err.Error())
	}
}
