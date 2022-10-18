package note

import (
	"neatly/internal/model/note"
	"neatly/pkg/logging"
)

type mapper struct {
	logger logging.Logger
}

func New(logger logging.Logger) *mapper {
	return &mapper{logger: logger}
}

func (m *mapper) MapCreateNoteDTO(dto note.CreateNoteDTO) note.Note {
	n := note.Note{
		ID:        0,
		Header:    dto.Header,
		Body:      dto.Body,
		ShortBody: "",
		Tags:      nil,
		Color:     dto.Color,
	}

	n.GenerateShortBody()
	m.logger.Infof("Generated short body with length of %v symbols", len(n.ShortBody))

	return n
}

func (m *mapper) MapGetAllNotesDTO(ns []note.Note) note.GetAllNotesDTO {
	return note.GetAllNotesDTO{
		Notes: ns,
	}
}

func (m *mapper) MapUpdateNoteDTO(dto note.UpdateNoteDTO) note.Note {
	if dto.Color == "" {
		dto.Color = note.DefaultNoteColor
	}

	n := note.Note{
		ID:        dto.ID,
		Header:    dto.Header,
		Body:      dto.Body,
		ShortBody: "",
		Color:     dto.Color,
	}

	n.GenerateShortBody()
	m.logger.Infof("Generated short body with length of %v symbols", len(n.ShortBody))

	return n
}