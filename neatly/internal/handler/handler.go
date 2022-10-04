package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "neatly/docs"
	"neatly/internal/mapper"
	"neatly/internal/service"
	"neatly/pkg/logging"

	"github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const (
	authURLGroup  = "/auth"
	registerURL   = "/register"
	loginURL      = "/login"
	apiURLGroup   = "/api"
	notesURLGroup = "/notes"
	tagsURLGroup  = "/tags"
	tagID         = "tag_id"
	searchURL     = "/search"
	versionAPI    = "1"
)

type Handler struct {
	logger   logging.Logger
	services *service.Service
	mappers  *mapper.Mapper
}

func New(services *service.Service, mappers *mapper.Mapper, logger logging.Logger) *Handler {
	return &Handler{
		services: services,
		logger:   logger,
		mappers:  mappers,
	}
}

func (h *Handler) RegisterHandler(idDebug *bool) *gin.Engine {
	if *idDebug == false {
		h.logger.Info("Setting gin to release mode")
		gin.SetMode(gin.ReleaseMode)
	}

	h.logger.Info("Create new gin router")
	router := gin.New()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group(fmt.Sprintf("%v/v%v", apiURLGroup, versionAPI))
	{
		auth := api.Group(authURLGroup)
		{
			auth.POST(registerURL, h.register)
			auth.POST(loginURL, h.login)
		}

		search := api.Group(searchURL, h.authenticate)
		{
			search.GET("", h.search)
		}

		notes := api.Group(notesURLGroup, h.authenticate)
		{
			notes.GET("", h.getAllNotes)       // /api/v1/notes
			notes.POST("", h.createNote)       // /api/v1/notes
			notes.GET("/:id", h.getOneNote)    // /api/v1/notes/:id
			notes.PATCH("/:id", h.updateNote)  // /api/v1/notes/:id
			notes.DELETE("/:id", h.deleteNote) // /api/v1/notes/:id

			tagsOnNote := notes.Group(fmt.Sprintf(":id%s", tagsURLGroup))
			{
				tagsOnNote.GET("", h.getAllTagsOnNote) // /api/notes/:id/tags/
				tagsOnNote.POST("", h.createTag)       // /api/notes/:id/tags/
			}
		}
		tags := api.Group(tagsURLGroup, h.authenticate)
		{
			tags.GET("", h.getAllTags)
			tags.GET("/:id", h.getOneTag)
			tags.PATCH("/:id", h.updateTag)
			tags.DELETE("/:id", h.deleteTag)
		}
	}

	return router
}
