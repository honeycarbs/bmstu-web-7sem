package mapper

import (
	"neatly/internal/model"
	"neatly/internal/model/dto"
	"neatly/pkg/logging"
)

type TagMapper struct {
	logger logging.Logger
}

func NewTagMapper(logger logging.Logger) *TagMapper {
	return &TagMapper{logger: logger}
}

func (m *TagMapper) MapCreateTagDTO(dto dto.CreateTagDTO) model.Tag {
	return model.Tag{
		ID:    0,
		Label: dto.Label,
		//Color: dto.Color,
	}
}

func (m *TagMapper) MapUpdateTagDTO(dto dto.UpdateTagDTO) model.Tag {
	return model.Tag{
		ID:    0,
		Label: dto.Label,
		//Color: dto.Color,
	}
}

func (m *TagMapper) MapGetAllTagsDTO(tags []model.Tag) dto.GetAllTagsDTO {
	return dto.GetAllTagsDTO{
		Tags: tags,
	}
}
