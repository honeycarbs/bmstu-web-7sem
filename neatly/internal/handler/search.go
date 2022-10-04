package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"neatly/pkg/e"
	"net/http"
)

const (
	tagSearchKey = "tag"
)

// @Summary Search Note
// @Security ApiKeyAuth
// @Tags search
// @Description search note
// @ID search-note
// @Accept  json
// @Produce json
// @Param   tag query  string  false  "notes search by tag"
// @Success 200 {object} note.GetAllNotesDTO
// @Failure 500 {object} e.ErrorResponse
// @Failure default {object} e.ErrorResponse
// @Router /api/v1/notes/search [get]
func (h *Handler) search(ctx *gin.Context) {
	userID, err := h.getUserID(ctx)
	if err != nil {
		h.logger.Info(err)
		return
	}

	keys := ctx.Request.URL.Query()
	values := keys[tagSearchKey]
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
