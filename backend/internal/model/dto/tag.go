package dto

import (
	"neatly/internal/model"
)

type CreateTagDTO struct {
	Label string `json:"label" binding:"required"`
}

type UpdateTagDTO struct {
	Label string `json:"label" db:"label"`
}

type GetAllTagsDTO struct {
	Tags []model.Tag `json:"tags"`
}
