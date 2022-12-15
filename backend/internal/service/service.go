package service

import (
	"neatly/internal/model"
	"neatly/internal/repository"
	"neatly/internal/service/account"
	"neatly/internal/service/note"
	"neatly/internal/service/tag"
	"neatly/pkg/logging"
)

type AccountService interface {
	CreateAccount(a *model.Account) error
	GenerateJWT(a *model.Account) (string, error)
}

type AccountServiceImpl struct {
	AccountService
}

func NewAccountServiceImpl(repo *repository.AccountRepositoryImpl, logger logging.Logger) *AccountServiceImpl {
	return &AccountServiceImpl{
		AccountService: account.NewService(repo, logger),
	}
}

type NoteService interface {
	Create(userID int, n *model.Note) error
	GetAll(userID int) ([]model.Note, error)
	GetOne(userID, noteID int) (model.Note, error)
	Delete(userID, noteID int) error
	Update(userID int, n model.Note, needBodyUpdate bool) error
	FindByTags(userID int, tagNames []string) ([]model.Note, error)
}

type NoteServiceImpl struct {
	NoteService
}

func NewNoteServiceImpl(noteRepo *repository.NoteRepositoryImpl, tagRepo *repository.TagRepositoryImpl, logger logging.Logger) *NoteServiceImpl {
	return &NoteServiceImpl{
		NoteService: note.NewService(noteRepo, tagRepo, logger),
	}
}

type TagService interface {
	Create(userID, noteID int, tag *model.Tag) (bool, error)
	GetAll(userID int) ([]model.Tag, error)
	GetAllByNote(userID, noteID int) ([]model.Tag, error)
	GetOne(userID, tagID int) (model.Tag, error)
	Delete(userID, tagID int) error
	Update(userID, tagID int, t model.Tag) error
	Detach(userID, tagID, noteID int) error
}

type TagServiceImpl struct {
	TagService
}

func NewTagServiceImpl(noteRepo *repository.NoteRepositoryImpl, tagRepo *repository.TagRepositoryImpl, logger logging.Logger) *TagServiceImpl {
	return &TagServiceImpl{
		TagService: tag.NewService(tagRepo, noteRepo, logger),
	}
}
