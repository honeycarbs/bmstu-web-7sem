package tag

import (
	"errors"
	"fmt"
	"neatly/internal/handlers/middleware"
	"neatly/internal/mapper"
	"neatly/internal/model/note"
	"neatly/internal/model/tag"
	"neatly/internal/service"
	"neatly/pkg/e"
	"neatly/pkg/logging"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	apiURLGroup   = "/api"
	notesURLGroup = "/notes"
	tagsURLGroup  = "/tags"
	apiVersion    = "1"
)

type Handler struct {
	logger  logging.Logger
	service service.Tag
	mapper  mapper.Tag
}

func NewHandler(logger logging.Logger, service service.Tag, mapper mapper.Tag) *Handler {
	return &Handler{logger: logger, service: service, mapper: mapper}
}

func (h *Handler) Register(router *gin.Engine) {
	tagsGroupName := fmt.Sprintf("%v/v%v%v", apiURLGroup, apiVersion, tagsURLGroup)
	tagsOnNoteGroupName := fmt.Sprintf("%v/v%v%v/:id%v", apiURLGroup, apiVersion, notesURLGroup, tagsURLGroup)

	h.logger.Tracef("Register route: %v", tagsGroupName)
	h.logger.Tracef("Register route: %v", tagsOnNoteGroupName)

	tagsGroup := router.Group(tagsGroupName, middleware.Authenticate)
	{
		tagsGroup.GET("", h.getAllTags)
		tagsGroup.GET("/:id", h.getOneTag)
		tagsGroup.PATCH("/:id", h.updateTag)
		tagsGroup.DELETE("/:id", h.deleteTag)
	}

	h.logger.Tracef("Register route: %v", tagsOnNoteGroupName)
	tagsOnNoteGroup := router.Group(tagsOnNoteGroupName, middleware.Authenticate)
	{
		tagsOnNoteGroup.GET("", h.getAllTagsOnNote)
		tagsOnNoteGroup.POST("", h.createTag)           // /api/notes/:id/tags/
		tagsOnNoteGroup.DELETE("/:tag_id", h.detachTag) // /api/notes/:id/tags/
	}
}

// @Summary Create tag
// @Security ApiKeyAuth
// @Tags tags
// @Description create tag
// @Accept  json
// @Produce  json
// @Param   id  path  string  true  "id"
// @Param dto body tag.CreateTagDTO true "tag info"
// @Success 201 {string} string 1
// @Failure 500 {object}  e.ErrorResponse
// @Failure 400,404 {object} e.ErrorResponse
// @Failure default {object}  e.ErrorResponse
// @Router /api/v1/notes/{id}/tags [post]
func (h *Handler) createTag(ctx *gin.Context) {
	userID, err := middleware.GetUserID(ctx)
	if err != nil {
		h.logger.Info(err)
		return
	}

	noteID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.logger.Info("error while getting id from request")
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	var (
		dto tag.CreateTagDTO
		t   tag.Tag
	)

	if err := ctx.BindJSON(&dto); err != nil {
		h.logger.Info(err)
		e.NewErrorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	t = h.mapper.MapCreateTagDTO(dto)
	err = h.service.Create(userID, noteID, &t)

	if err != nil {
		if errors.Is(err, &note.NoteNotFoundErr{}) {
			e.NewErrorResponse(ctx, http.StatusNotFound, err)
			return
		}
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusCreated, fmt.Sprintf(
		"%s%s/%v", apiURLGroup, tagsURLGroup, t.ID))
}

// @Summary Get all tags on one note
// @Security ApiKeyAuth
// @Tags tags
// @Description get tags for note
// @Accept  json
// @Produce  json
// @Param   id  path  string  true  "id"
// @Success 200 {object} tag.GetAllTagsDTO
// @Failure 500 {object}  e.ErrorResponse
// @Failure 400,404 {object} e.ErrorResponse
// @Failure default {object}  e.ErrorResponse
// @Router /api/v1/notes/{id}/tags [get]
func (h *Handler) getAllTagsOnNote(ctx *gin.Context) {
	userID, err := middleware.GetUserID(ctx)
	if err != nil {
		h.logger.Info(err)
		return
	}

	noteID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.logger.Info("error while getting id from request")
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	tags, err := h.service.GetAllByNote(userID, noteID)

	if err != nil {
		h.logger.Info(err)
		if errors.Is(err, &note.NoteNotFoundErr{}) {
			e.NewErrorResponse(ctx, http.StatusNotFound, err)
			return
		}
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	dto := h.mapper.MapGetAllTagsDTO(tags)

	ctx.JSON(http.StatusOK, dto)
}

// @Summary Get all tags
// @Security ApiKeyAuth
// @Tags tags
// @Description get tags from user
// @Accept  json
// @Produce  json
// @Success 200 {object} tag.GetAllTagsDTO
// @Failure 500 {object}  e.ErrorResponse
// @Failure 400,404 {object} e.ErrorResponse
// @Failure default {object}  e.ErrorResponse
// @Router /api/v1/tags [get]
func (h *Handler) getAllTags(ctx *gin.Context) {
	userID, err := middleware.GetUserID(ctx)
	if err != nil {
		h.logger.Info(err)
		return
	}

	tags, err := h.service.GetAll(userID)

	if err != nil {
		h.logger.Info(err)
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	dto := h.mapper.MapGetAllTagsDTO(tags)

	ctx.JSON(http.StatusOK, dto)
}

// @Summary Get one tag by ID
// @Security ApiKeyAuth
// @Tags tags
// @Description get one tag by ID
// @Accept  json
// @Produce  json
// @Param   id  path  string  true  "id"
// @Success 200 {object} tag.Tag
// @Failure 500 {object}  e.ErrorResponse
// @Failure 400,404 {object} e.ErrorResponse
// @Failure default {object}  e.ErrorResponse
// @Router /api/v1/tags/{id} [get]
func (h *Handler) getOneTag(ctx *gin.Context) {
	userID, err := middleware.GetUserID(ctx)
	if err != nil {
		h.logger.Info(err)
		return
	}

	tagID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.logger.Info("error while getting id from request")
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	t, err := h.service.GetOne(userID, tagID)

	if err != nil {
		h.logger.Info(err)
		if errors.Is(err, &tag.TagNotFoundErr{}) {
			e.NewErrorResponse(ctx, http.StatusNotFound, err)
			return
		}
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, t)
}

// @Summary Update tag by ID
// @Security ApiKeyAuth
// @Tags tags
// @Description update one tag by ID
// @Accept  json
// @Produce  json
// @Param   id  path  string  true  "id"
// @Param dto body tag.UpdateTagDTO true "tag info"
// @Success 204
// @Failure 500 {object}  e.ErrorResponse
// @Failure 400,404 {object} e.ErrorResponse
// @Failure default {object}  e.ErrorResponse
// @Router /api/v1/tags/{id} [patch]
func (h *Handler) updateTag(ctx *gin.Context) {
	userID, err := middleware.GetUserID(ctx)
	if err != nil {
		return
	}

	tagID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.logger.Info("error while getting id from request")
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	var (
		dto tag.UpdateTagDTO
	)
	if err := ctx.BindJSON(&dto); err != nil {
		h.logger.Info(err)
		e.NewErrorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	t := h.mapper.MapUpdateTagDTO(dto)
	err = h.service.Update(userID, tagID, t)
	if err != nil {
		h.logger.Info(err)
		if errors.Is(err, &tag.TagNotFoundErr{}) {
			e.NewErrorResponse(ctx, http.StatusNotFound, err)
			return
		}
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Writer.WriteHeader(http.StatusNoContent)
}

// @Summary Delete one tag by ID
// @Security ApiKeyAuth
// @Tags tags
// @Description delete one tag by ID
// @Accept  json
// @Produce  json
// @Param   id  path  string  true  "id"
// @Success 200 {integer} integer 1
// @Failure 500 {object}  e.ErrorResponse
// @Failure 400,404 {object} e.ErrorResponse
// @Failure default {object}  e.ErrorResponse
// @Router /api/v1/tags/{id} [delete]
func (h *Handler) deleteTag(ctx *gin.Context) {
	userID, err := middleware.GetUserID(ctx)
	if err != nil {
		h.logger.Info(err)
		return
	}

	tagID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.logger.Info("error while getting id from request")
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	err = h.service.Delete(userID, tagID)

	if err != nil {
		h.logger.Info(err)
		if errors.Is(err, &tag.TagNotFoundErr{}) {
			e.NewErrorResponse(ctx, http.StatusNotFound, err)
			return
		}
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Writer.WriteHeader(http.StatusNoContent)
}

// @Summary Detach tag by ID from note by ID
// @Security ApiKeyAuth
// @Tags tags
// @Description detach tag by ID from note by ID
// @Accept  json
// @Produce  json
// @Param   id  path  string  true  "id"
// @Param   tag_id  path  string  true  "tag id"
// @Success 200 {integer} integer 1
// @Failure 500 {object}  e.ErrorResponse
// @Failure 400,404 {object} e.ErrorResponse
// @Failure default {object}  e.ErrorResponse
// @Router /api/v1/tags/{id}/tags/{tag_id} [delete]
func (h *Handler) detachTag(ctx *gin.Context) {
	userID, err := middleware.GetUserID(ctx)
	if err != nil {
		h.logger.Info(err)
		return
	}

	noteID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.logger.Info("error while getting id from request")
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	tagID, err := strconv.Atoi(ctx.Param("tag_id"))
	if err != nil {
		h.logger.Info("error while getting id from request")
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	err = h.service.Detach(userID, tagID, noteID)

	if err != nil {
		h.logger.Info(err)
		if errors.Is(err, &tag.TagNotFoundErr{}) {
			e.NewErrorResponse(ctx, http.StatusNotFound, err)
			return
		}
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Writer.WriteHeader(http.StatusNoContent)

}
