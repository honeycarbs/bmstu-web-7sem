package tag

import (
	"neatly/internal/model/tag"
	"neatly/pkg/logging"
)

type mapper struct {
	logger logging.Logger
}

func New(logger logging.Logger) *mapper {
	return &mapper{logger: logger}
}

func (m *mapper) MapCreateTagDTO(dto tag.CreateTagDTO) tag.Tag {
	return tag.Tag{
		ID:    0,
		Name:  dto.Name,
		Color: dto.Color,
	}
}

func (m *mapper) MapUpdateTagDTO(dto tag.UpdateTagDTO) tag.Tag {
	return tag.Tag{
		ID:    0,
		Name:  dto.Name,
		Color: dto.Color,
	}
}
