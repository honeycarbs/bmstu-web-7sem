package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"neatly/internal/model/note"
	"neatly/pkg/e"
	"net/http"
	"strconv"
)

// @Summary Create note
// @Security ApiKeyAuth
// @Tags notes
// @Description create note
// @Accept  json
// @Produce  json
// @Param dto body note.CreateNoteDTO true "note content"
// @Success 201 {string} string 1
// @Failure 500 {object}  e.ErrorResponse
// @Failure 400,404 {object} e.ErrorResponse
// @Failure default {object}  e.ErrorResponse
// @Router /api/v1/notes [post]
func (h *Handler) createNote(ctx *gin.Context) {
	userID, err := h.getUserID(ctx)
	if err != nil {
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err)
		h.logger.Info(err)
		return
	}

	var dto note.CreateNoteDTO
	if err := ctx.BindJSON(&dto); err != nil {
		h.logger.Info(err)
		if errors.Is(err, &e.NoteNotFoundErr{}) {
			e.NewErrorResponse(ctx, http.StatusNotFound, err)
			return
		}
		e.NewErrorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	n := h.mappers.Note.MapCreateNoteDTO(dto)
	err = h.services.Note.Create(userID, &n)
	if err != nil {
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusCreated, fmt.Sprintf(
		"%s%s/%v", apiURLGroup, notesURLGroup, n.ID))
}

// @Summary Get all notes from user
// @Security ApiKeyAuth
// @Tags notes
// @Description create note
// @Accept  json
// @Produce  json
// @Success 201 {string} string 1
// @Failure 500 {object}  e.ErrorResponse
// @Failure 400,404 {object} e.ErrorResponse
// @Failure default {object}  e.ErrorResponse
// @Router  /api/v1/notes [get]
func (h *Handler) getAllNotes(ctx *gin.Context) {
	userID, err := h.getUserID(ctx)
	if err != nil {
		if errors.Is(err, &e.AccountNotFoundErr{}) {
			e.NewErrorResponse(ctx, http.StatusNotFound, err)
			return
		}
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	notes, err := h.services.Note.GetAll(userID)
	if err != nil {
		if errors.Is(err, &e.NoteNotFoundErr{}) {
			e.NewErrorResponse(ctx, http.StatusNotFound, err)
			return
		}
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}
	ndto := h.mappers.Note.MapGetAllNotesDTO(notes)

	ctx.JSON(http.StatusOK, ndto)
}

// @Summary Get Note By Id
// @Security ApiKeyAuth
// @Tags notes
// @Description get note by id
// @ID get-note-by-id
// @Accept  json
// @Produce json
// @Param   id  path  string  true  "id"
// @Success 200 {object} note.Note
// @Failure 500 {object} e.ErrorResponse
// @Failure default {object} e.ErrorResponse
// @Router /api/v1/notes/{id} [get]
func (h *Handler) getOneNote(ctx *gin.Context) {
	userID, err := h.getUserID(ctx)
	if err != nil {
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	noteID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.logger.Info("error while getting id from request")
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	n, err := h.services.Note.GetOne(userID, noteID)
	if err != nil {
		if errors.Is(err, &e.NoteNotFoundErr{}) {
			e.NewErrorResponse(ctx, http.StatusNotFound, err)
			return
		}
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, n)
}

// @Summary Update Note
// @Security ApiKeyAuth
// @Tags notes
// @Description update note
// @ID update-note
// @Accept  json
// @Produce json
// @Param   id   path  string  true  "id"
// @Param dto body note.UpdateNoteDTO true "note content"
// @Success 204
// @Failure 500 {object} e.ErrorResponse
// @Failure default {object} e.ErrorResponse
// @Router /api/v1/notes/{id} [patch]
func (h *Handler) updateNote(ctx *gin.Context) {
	userID, err := h.getUserID(ctx)
	if err != nil {
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	noteID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.logger.Info("error while getting id from request")
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	bodyBytes, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		h.logger.Info(err)
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err)
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
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	if err := json.Unmarshal(bodyBytes, &data); err != nil {
		h.logger.Info(err)
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	_, needBodyUpdate = data["body"]
	h.logger.Infof("Need body update: %v", needBodyUpdate)

	n := h.mappers.Note.MapUpdateNoteDTO(dto)
	err = h.services.Note.Update(userID, n, needBodyUpdate)

	if err != nil {
		h.logger.Info(err)
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Writer.WriteHeader(http.StatusNoContent)
}

// @Summary Delete Note
// @Security ApiKeyAuth
// @Tags notes
// @Description delete note
// @ID delete-note
// @Accept  json
// @Produce json
// @Param   id   path string  true  "id"
// @Success 200 {integer} integer 1
// @Failure 500 {object} e.ErrorResponse
// @Failure default {object} e.ErrorResponse
// @Router /api/v1/notes/{id} [delete]
func (h *Handler) deleteNote(ctx *gin.Context) {
	userID, err := h.getUserID(ctx)
	if err != nil {
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	noteID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.logger.Info("error while getting id from request")
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	err = h.services.Note.Delete(userID, noteID)
	if err != nil {
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, noteID)
}
