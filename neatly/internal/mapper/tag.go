package mapper

import (
	"neatly/internal/model/tag"
	"neatly/pkg/logging"
)

type TagMapper struct {
	logger logging.Logger
}

func NewTagMapper(logger logging.Logger) *TagMapper {
	return &TagMapper{logger: logger}
}

func (m *TagMapper) MapCreateTagDTO(dto tag.CreateTagDTO) tag.Tag {
	return tag.Tag{
		ID:    0,
		Name:  dto.Name,
		Color: dto.Color,
	}
}

func (m *TagMapper) MapUpdateTagDTO(dto tag.UpdateTagDTO) tag.Tag {
	return tag.Tag{
		ID:    0,
		Name:  dto.Name,
		Color: dto.Color,
	}
}
