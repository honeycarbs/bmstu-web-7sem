package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"neatly"
	_ "neatly/docs"
	"neatly/internal/handlers/account"
	"neatly/internal/handlers/note"
	"neatly/internal/handlers/tag"
	"neatly/internal/mapper"
	"neatly/internal/repository"
	"neatly/internal/service"
	"neatly/internal/session"
	"neatly/pkg/client/psqlclient"
	"neatly/pkg/logging"
)

// @title Neat.ly API
// @version 1.0
// @description API Server for notes-taking applications
// @termsOfService  http://swagger.io/terms/

// @host localhost:8000
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	logging.Init()
	logger := logging.GetLogger()

	cfg := session.GetConfig()

	client, err := psqlclient.NewClient(cfg.DB)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Info("Create new gin router")
	router := gin.New()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	logger.Info("initializing repository")
	repos := repository.New(client, logger)

	logger.Info("initializing services")
	services := service.New(repos, logger)
	mappers := mapper.New(logger)

	accountHandler := account.NewHandler(logger, services.Account, mappers.Account)
	accountHandler.Register(router)

	notesHandler := note.NewHandler(logger, services.Note, mappers.Note)
	notesHandler.Register(router)

	tagsHandler := tag.NewHandler(logger, services.Tag, mappers.Tag)
	tagsHandler.Register(router)

	backend.Run(cfg, router, logger)
}
