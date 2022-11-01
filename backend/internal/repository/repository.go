package repository

import (
	"neatly/internal/model"
	"neatly/internal/repository/psql"
	"neatly/pkg/client/psqlclient"
	"neatly/pkg/logging"
)

//go:generate mockgen -destination=mock/$GOFILE -package=mock -source=$GOFILE

type AccountRepository interface {
	CreateAccount(a *model.Account) error
	AuthorizeAccount(a *model.Account) error
	GetOne(userID int) (model.Account, error)
}

type AccountRepositoryImpl struct {
	AccountRepository
}

func NewAccountRepositoryImpl(client *psqlclient.Client, logger logging.Logger) *AccountRepositoryImpl {
	return &AccountRepositoryImpl{
		AccountRepository: psql.NewAccountPostgres(client, logger),
	}
}

type NoteRepository interface {
	Create(userID int, note *model.Note) error
	GetAll(userID int) ([]model.Note, error)
	GetOne(userID, noteID int) (model.Note, error)
	Delete(userID, noteID int) error
	Update(userID int, n model.Note) error
}

type NoteRepositoryImpl struct {
	NoteRepository
}

func NewNoteRepositoryImpl(client *psqlclient.Client, logger logging.Logger) *NoteRepositoryImpl {
	return &NoteRepositoryImpl{NoteRepository: psql.NewNotePostgres(client, logger)}
}

type TagRepository interface {
	Create(userID int, noteID int, t *model.Tag) error
	GetAll(userID int) ([]model.Tag, error)
	GetAllByNote(userID, noteID int) ([]model.Tag, error)
	GetOne(userID, tagID int) (model.Tag, error)
	Delete(userID, tagID int) error
	Detach(userID, tagID, noteID int) error
	Assign(tagID, noteID, userID int) error
	Update(userID, tagID int, t model.Tag) error
}

type TagRepositoryImpl struct {
	TagRepository
}

func NewTagRepositoryImpl(client *psqlclient.Client, logger logging.Logger) *TagRepositoryImpl {
	return &TagRepositoryImpl{
		TagRepository: psql.NewTagPostgres(client, logger),
	}
}
