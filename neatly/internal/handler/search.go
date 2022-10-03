package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"neatly/internal/e"
	"net/http"
)

const (
	tagSearchKey    = "tag"
	headerSearchKey = "header"
)

func (h *Handler) search(ctx *gin.Context) {
	userID, err := h.getUserID(ctx)
	if err != nil {
		h.logger.Info(err)
		return
	}

	keys := ctx.Request.URL.Query()
	values := keys["tag"]
	if values == nil {
		ctx.JSON(http.StatusBadRequest, fmt.Sprint("malformed query"))
	}

	h.logger.Infof("Values got from req: %v", values)
	ns, err := h.services.Note.FindByTags(userID, values)
	if err != nil {
		h.logger.Info(err)
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}
	dto := h.mappers.Note.MapGetAllNotesDTO(ns)

	ctx.JSON(http.StatusOK, dto)
}
