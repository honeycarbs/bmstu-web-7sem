package note

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"neatly/internal/handlers/middleware"
	"neatly/internal/mapper"
	"neatly/internal/model/account"
	"neatly/internal/model/note"
	"neatly/internal/service"
	"neatly/pkg/e"
	"neatly/pkg/logging"
	"net/http"
	"strconv"
)

const (
	tagsURLGroup  = "/tags"
	notesURLGroup = "/notes"
	apiURLGroup   = "/api"
	apiVersion    = "1"
	searchURL     = "/search"
	tagSearchKey  = "tag"
)

type Handler struct {
	logger  logging.Logger
	service service.Note
	mapper  mapper.Note
}

func NewHandler(logger logging.Logger, service service.Note, mapper mapper.Note) *Handler {
	return &Handler{logger: logger, service: service, mapper: mapper}
}

func (h *Handler) Register(router *gin.Engine) {
	groupName := fmt.Sprintf("%v/v%v/%v", apiURLGroup, apiVersion, notesURLGroup)

	h.logger.Tracef("Register route: %v", groupName)

	group := router.Group(groupName, middleware.Authenticate)
	{
		group.GET(searchURL, h.search)
		group.GET("", h.getAllNotes)       // /api/v1/notes
		group.POST("", h.createNote)       // /api/v1/notes
		group.GET("/:id", h.getOneNote)    // /api/v1/notes/:id
		group.PATCH("/:id", h.updateNote)  // /api/v1/notes/:id
		group.DELETE("/:id", h.deleteNote) // /api/v1/notes/:id
	}
}

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
	userID, err := middleware.GetUserID(ctx)
	if err != nil {
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err)
		h.logger.Info(err)
		return
	}

	var dto note.CreateNoteDTO
	if err := ctx.BindJSON(&dto); err != nil {
		h.logger.Info(err)
		if errors.Is(err, &note.NoteNotFoundErr{}) {
			e.NewErrorResponse(ctx, http.StatusNotFound, err)
			return
		}
		e.NewErrorResponse(ctx, http.StatusBadRequest, err)
		return
	}

	n := h.mapper.MapCreateNoteDTO(dto)
	err = h.service.Create(userID, &n)
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
// @Success 200 {object} note.GetAllNotesDTO
// @Failure 500 {object}  e.ErrorResponse
// @Failure 400,404 {object} e.ErrorResponse
// @Failure default {object}  e.ErrorResponse
// @Router  /api/v1/notes [get]
func (h *Handler) getAllNotes(ctx *gin.Context) {
	userID, err := middleware.GetUserID(ctx)
	if err != nil {
		if errors.Is(err, &account.AccountNotFoundErr{}) {
			e.NewErrorResponse(ctx, http.StatusNotFound, err)
			return
		}
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	notes, err := h.service.GetAll(userID)
	if err != nil {
		if errors.Is(err, &note.NoteNotFoundErr{}) {
			e.NewErrorResponse(ctx, http.StatusNotFound, err)
			return
		}
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}
	ndto := h.mapper.MapGetAllNotesDTO(notes)

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
	userID, err := middleware.GetUserID(ctx)
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

	n, err := h.service.GetOne(userID, noteID)
	if err != nil {
		if errors.Is(err, &note.NoteNotFoundErr{}) {
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
	userID, err := middleware.GetUserID(ctx)
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

	n := h.mapper.MapUpdateNoteDTO(dto)
	err = h.service.Update(userID, n, needBodyUpdate)

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
// @Success 204
// @Failure 500 {object} e.ErrorResponse
// @Failure default {object} e.ErrorResponse
// @Router /api/v1/notes/{id} [delete]
func (h *Handler) deleteNote(ctx *gin.Context) {
	userID, err := middleware.GetUserID(ctx)
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

	err = h.service.Delete(userID, noteID)
	if err != nil {
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Writer.WriteHeader(http.StatusNoContent)
}

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
	userID, err := middleware.GetUserID(ctx)
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
	ns, err := h.service.FindByTags(userID, values)
	if err != nil {
		h.logger.Info(err)
		e.NewErrorResponse(ctx, http.StatusInternalServerError, err)
		return
	}
	dto := h.mapper.MapGetAllNotesDTO(ns)

	ctx.JSON(http.StatusOK, dto)
}
