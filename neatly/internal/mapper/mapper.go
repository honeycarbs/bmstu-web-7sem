package mapper

import (
	authMapper "neatly/internal/mapper/auth"
	noteMapper "neatly/internal/mapper/note"
	tagMapper "neatly/internal/mapper/tag"
	"neatly/internal/model/auth"
	"neatly/internal/model/note"
	"neatly/internal/model/tag"
	"neatly/pkg/logging"
)

type Auth interface {
	MapRegisterAccountDTO(dto auth.RegisterAccountDTO) (auth.Account, error)
	MapLogInAccountDTO(dto auth.LoginAccountDTO) auth.Account
	MapJwtDTO(token string) auth.JwtDTO
}

type Note interface {
	MapCreateNoteDTO(dto note.CreateNoteDTO) note.Note
	MapUpdateNoteDTO(dto note.UpdateNoteDTO) note.Note
	MapGetAllNotesDTO(ns []note.Note) note.GetAllNotesDTO
}

type Tag interface {
	MapCreateTagDTO(dto tag.CreateTagDTO) tag.Tag
	MapUpdateTagDTO(dto tag.UpdateTagDTO) tag.Tag
	MapGetAllTagsDTO(tags []tag.Tag) tag.GetAllTagsDTO
}

type Mapper struct {
	Auth
	Note
	Tag
}

func New(l logging.Logger) *Mapper {
	return &Mapper{
		Auth: authMapper.New(l),
		Note: noteMapper.New(l),
		Tag:  tagMapper.New(l),
	}
}
