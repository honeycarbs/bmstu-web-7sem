package dto

import (
	"neatly/internal/model"
)

type CreateNoteDTO struct {
	Header string `json:"header" binding:"required"`
	Color  string `json:"color" binding:"required"`
	Body   string `json:"body"`
}

type UpdateNoteDTO struct {
	ID     int
	Header string `json:"header"`
	Body   string `json:"body"`
	Color  string `json:"color"`
}

type GetAllNotesDTO struct {
	Notes []model.Note `json:"notes"`
}
