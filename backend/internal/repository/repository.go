package repository

import (
	"neatly/internal/model/account"
	"neatly/internal/model/note"
	"neatly/internal/model/tag"
	"neatly/internal/repository/psql"
	"neatly/pkg/client/psqlclient"
	"neatly/pkg/logging"
)

type Account interface {
	CreateAccount(a *account.Account) error
	GetAccount(a *account.Account) error
}

type Note interface {
	Create(userID int, note *note.Note) error
	GetAll(userID int) ([]note.Note, error)
	GetOne(userID, noteID int) (note.Note, error)
	Delete(userID, noteID int) error
	Update(userID int, n note.Note) error
}

type Tag interface {
	Create(userID int, noteID int, t *tag.Tag) error
	GetAll(userID int) ([]tag.Tag, error)
	GetAllByNote(userID, noteID int) ([]tag.Tag, error)
	GetOne(userID, tagID int) (tag.Tag, error)
	Delete(userID, tagID int) error
	Detach(userID, tagID, noteID int) error
	Assign(tagID, noteID, userID int) error
	Update(userID, tagID int, t tag.Tag) error
}

type Repository struct {
	Account
	Note
	Tag
}

func New(client *psqlclient.Client, logger logging.Logger) *Repository {
	return &Repository{
		Account: psql.NewAuthPostgres(client, logger),
		Note:    psql.NewNotePostgres(client, logger),
		Tag:     psql.NewTagPostgres(client, logger),
	}
}
