package service

import (
	"neatly/internal/model/note"
	"neatly/internal/model/tag"
	"neatly/internal/model/user"
	"neatly/internal/repository"
	"neatly/internal/service/auth"
	noteService "neatly/internal/service/note"
	tagService "neatly/internal/service/tag"
	"neatly/pkg/logging"
)

type Authorisation interface {
	CreateUser(u *user.User) error
	GenerateJWT(u *user.User) (string, error)
}

type Note interface {
	Create(userID int, n *note.Note) error
	GetAll(userID int) ([]note.Note, error)
	GetOne(userID, noteID int) (note.Note, error)
	Delete(userID, noteID int) error
	Update(userID int, n note.Note, needBodyUpdate bool) error
	FindByTags(userID int, tagNames []string) ([]note.Note, error)
}

type Tag interface {
	Create(userID, noteID int, tag *tag.Tag) error
	GetAll(userID int) ([]tag.Tag, error)
	GetAllByNote(userID, noteID int) ([]tag.Tag, error)
	GetOne(userID, tagID int) (tag.Tag, error)
	Delete(userID, tagID int) error
	Update(userID, tagID int, t tag.Tag) error
}

type Service struct {
	Authorisation
	Note
	Tag
}

func New(repo *repository.Repository, logger logging.Logger) *Service {
	return &Service{
		Authorisation: auth.NewService(repo.Authorisation),
		Note:          noteService.NewService(repo.Note, repo.Tag, logger),
		Tag:           tagService.NewService(repo.Tag, repo.Note, logger),
	}
}
