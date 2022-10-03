package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"neatly/internal/e"
	"neatly/internal/model/tag"
	"net/http"
	"strconv"
)

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
