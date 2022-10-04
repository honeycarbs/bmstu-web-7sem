package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"neatly/internal/e"
	"neatly/internal/model/tag"
	"net/http"
	"strconv"
)

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
	userID, err := h.getUserID(ctx)
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

	t = h.mappers.Tag.MapCreateTagDTO(dto)
	err = h.services.Tag.Create(userID, noteID, &t)

	if err != nil {
		h.logger.Info(err)
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
	userID, err := h.getUserID(ctx)
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

	tags, err := h.services.Tag.GetAllByNote(userID, noteID)

	if err != nil {
		h.logger.Info(err)
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	dto := h.mappers.Tag.MapGetAllTagsDTO(tags)

	ctx.JSON(http.StatusCreated, dto)
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
	userID, err := h.getUserID(ctx)
	if err != nil {
		h.logger.Info(err)
		return
	}

	tags, err := h.services.Tag.GetAll(userID)

	if err != nil {
		h.logger.Info(err)
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	dto := h.mappers.Tag.MapGetAllTagsDTO(tags)

	ctx.JSON(http.StatusCreated, dto)
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
	userID, err := h.getUserID(ctx)
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

	tag, err := h.services.Tag.GetOne(userID, tagID)

	if err != nil {
		h.logger.Info(err)
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusCreated, tag)
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
	userID, err := h.getUserID(ctx)
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

	t := h.mappers.Tag.MapUpdateTagDTO(dto)
	err = h.services.Tag.Update(userID, tagID, t)
	if err != nil {
		h.logger.Info(err)
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
	userID, err := h.getUserID(ctx)
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

	err = h.services.Tag.Delete(userID, tagID)

	if err != nil {
		h.logger.Info(err)
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, tagID)
}
