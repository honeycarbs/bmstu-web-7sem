package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"neatly/cmd/server"
	"neatly/docs"
	_ "neatly/docs"
	"neatly/internal/handlers/account"
	"neatly/internal/handlers/middleware"
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

	logger.Info("Configure CORS")
	middleware.CorsMiddleware(router)

	docs.SwaggerInfo.Host = cfg.Swagger.Host
	router.GET("api/v1/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	logger.Info("initializing account repository")
	accountRepo := repository.NewAccountRepositoryImpl(client, logger)
	logger.Info("initializing note repository")
	noteRepo := repository.NewNoteRepositoryImpl(client, logger)
	logger.Info("initializing tag repository")
	tagRepo := repository.NewTagRepositoryImpl(client, logger)

	logger.Info("initializing account service")
	accountService := service.NewAccountServiceImpl(accountRepo, logger)
	logger.Info("initializing note service")
	noteService := service.NewNoteServiceImpl(noteRepo, tagRepo, logger)
	logger.Info("initializing tag service")
	tagService := service.NewTagServiceImpl(noteRepo, tagRepo, logger)

	logger.Info("initializing account mapper")
	accountMapper := mapper.NewAccountMapper(logger)
	logger.Info("initializing note mapper")
	noteMapper := mapper.NewNoteMapper(logger)
	logger.Info("initializing tag mapper")
	tagMapper := mapper.NewTagMapper(logger)

	logger.Info("initializing account handler")
	accountHandler := account.NewHandler(logger, *accountService, *accountMapper)
	accountHandler.Register(router)

	logger.Info("initializing note handler")
	noteHandler := note.NewHandler(logger, *noteService, *noteMapper)
	noteHandler.Register(router)

	logger.Info("initializing tag handler")
	tagHandler := tag.NewHandler(logger, tagService, *tagMapper)
	tagHandler.Register(router)

	server.Run(cfg, router, logger)
}
