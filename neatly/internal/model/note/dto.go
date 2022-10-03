package note

import (
	"neatly/internal/model/tag"
)

type CreateNoteDTO struct {
	Header string `json:"header" binding:"required"`
	Color  string `json:"color" binding:"required"`
	Body   string `json:"body"`
}

type UpdateNoteDTO struct {
	ID     int
	Header string    `json:"header"`
	Body   string    `json:"body"`
	Color  string    `json:"color"`
	Tags   []tag.Tag `json:"tags"`
}

type GetAllNotesDTO struct {
	Notes []Note `json:"notes"`
}
