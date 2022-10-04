package main

import (
	_ "github.com/lib/pq"
	"neatly"
	"neatly/internal/handler"
	"neatly/internal/mapper"
	"neatly/internal/repository"
	"neatly/internal/service"
	"neatly/internal/session"
	"neatly/pkg/logging"
	"neatly/pkg/postgres"
)

// @title Neat.ly API
// @version 1.0
// @description API Server for notes-taking applications
// @termsOfService  http://swagger.io/terms/

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
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

	logger.Info("initializing mappers")
	mappers := mapper.New(logger)

	logger.Info("initializing handler")
	handlers := handler.New(services, mappers, logger)

	neatly.Run(cfg, handlers.RegisterHandler(cfg.IsDebug), logger)
}
