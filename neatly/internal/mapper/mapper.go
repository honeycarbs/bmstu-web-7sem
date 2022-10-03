package mapper

import (
	noteMapper "neatly/internal/mapper/note"
	tagMapper "neatly/internal/mapper/tag"
	userMapper "neatly/internal/mapper/user"
	"neatly/internal/model/note"
	"neatly/internal/model/tag"
	"neatly/internal/model/user"
	"neatly/pkg/logging"
)

type User interface {
	MapUserRegisterDTO(dto user.RegisterUserDTO) (user.User, error)
	MapUserLogInUserDTO(dto user.LoginUserDTO) user.User
}

type Note interface {
	MapCreateNoteDTO(dto note.CreateNoteDTO) note.Note
	MapUpdateNoteDTO(dto note.UpdateNoteDTO) note.Note
}

type Tag interface {
	MapCreateTagDTO(dto tag.CreateTagDTO) tag.Tag
	MapUpdateTagDTO(dto tag.UpdateTagDTO) tag.Tag
}

type Mapper struct {
	User
	Note
	Tag
}

func NewMapper(l logging.Logger) *Mapper {
	return &Mapper{
		User: userMapper.New(l),
		Note: noteMapper.New(l),
		Tag:  tagMapper.New(l),
	}
}
