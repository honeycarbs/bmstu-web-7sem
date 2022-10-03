package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"neatly/internal/e"
	"neatly/internal/model/note"
	"net/http"
	"strconv"
)

func (h *Handler) createNote(ctx *gin.Context) {
	id, err := h.getUserID(ctx)
	if err != nil {
		h.logger.Info(err)
		return
	}

	var dto note.CreateNoteDTO
	if err := ctx.BindJSON(&dto); err != nil {
		h.logger.Info(err)
		e.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	n := h.mappers.Note.MapCreateNoteDTO(dto)
	err = h.services.Note.Create(id, &n)
	if err != nil {
		h.logger.Info(err)
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, fmt.Sprintf(
		"%s%s/%v", apiURLGroup, notesURLGroup, n.ID))
}

func (h *Handler) getAllNotes(ctx *gin.Context) {
	id, err := h.getUserID(ctx)
	if err != nil {
		return
	}

	notes, err := h.services.Note.GetAll(id)
	if err != nil {
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, notes)
}

func (h *Handler) getOneNote(ctx *gin.Context) {
	userID, err := h.getUserID(ctx)
	if err != nil {
		return
	}

	noteID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.logger.Info("error while getting id from request")
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	n, err := h.services.Note.GetOne(userID, noteID)
	if err != nil {
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, n)
}

func (h *Handler) updateNote(ctx *gin.Context) {
	userID, err := h.getUserID(ctx)
	if err != nil {
		return
	}

	noteID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.logger.Info("error while getting id from request")
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	bodyBytes, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		h.logger.Info(err)
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	h.logger.Debug("unmarshal body bytes")
	var (
		dto            note.UpdateNoteDTO
		data           map[string]interface{}
		needBodyUpdate bool
	)
	h.logger.Infof("NOTE ID: %v", noteID)
	dto.ID = noteID
	if err := json.Unmarshal(bodyBytes, &dto); err != nil {
		h.logger.Info(err)
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	// TODO: check for updating tags
	if err := json.Unmarshal(bodyBytes, &data); err != nil {
		h.logger.Info(err)
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	_, needBodyUpdate = data["body"]
	h.logger.Infof("Need body update: %v", needBodyUpdate)

	n := h.mappers.Note.MapUpdateNoteDTO(dto)
	err = h.services.Note.Update(userID, n, needBodyUpdate)

	if err != nil {
		h.logger.Info(err)
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.Writer.WriteHeader(http.StatusNoContent)
}

func (h *Handler) deleteNote(ctx *gin.Context) {
	userID, err := h.getUserID(ctx)
	if err != nil {
		return
	}

	noteID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.logger.Info("error while getting id from request")
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	err = h.services.Note.Delete(userID, noteID)
	if err != nil {
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, noteID)
}
