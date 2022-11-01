package mapper

import (
	"neatly/internal/model"
	"neatly/internal/model/dto"
	"neatly/pkg/logging"
)

type NoteMapper struct {
	logger logging.Logger
}

func NewNoteMapper(logger logging.Logger) *NoteMapper {
	return &NoteMapper{logger: logger}
}

func (m *NoteMapper) MapCreateNoteDTO(dto dto.CreateNoteDTO) model.Note {
	n := model.Note{
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

func (m *NoteMapper) MapGetAllNotesDTO(ns []model.Note) dto.GetAllNotesDTO {
	return dto.GetAllNotesDTO{
		Notes: ns,
	}
}

func (m *NoteMapper) MapUpdateNoteDTO(dto dto.UpdateNoteDTO) model.Note {
	if dto.Color == "" {
		dto.Color = model.DefaultNoteColor
	}

	n := model.Note{
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
