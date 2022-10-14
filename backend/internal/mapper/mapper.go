package mapper

import (
	authMapper "neatly/internal/mapper/account"
	noteMapper "neatly/internal/mapper/note"
	tagMapper "neatly/internal/mapper/tag"
	"neatly/internal/model/account"
	"neatly/internal/model/note"
	"neatly/internal/model/tag"
	"neatly/pkg/logging"
)

type Account interface {
	MapRegisterAccountDTO(dto account.RegisterAccountDTO) (account.Account, error)
	MapLogInAccountDTO(dto account.LoginAccountDTO) account.Account
	MapAccountWithTokenDTO(token string, a account.Account) account.WithTokenDTO
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
	Account
	Note
	Tag
}

func New(l logging.Logger) *Mapper {
	return &Mapper{
		Account: authMapper.New(l),
		Note:    noteMapper.New(l),
		Tag:     tagMapper.New(l),
	}
}
