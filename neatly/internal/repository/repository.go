package repository

import (
	"github.com/jmoiron/sqlx"
	"neatly/internal/model/auth"
	"neatly/internal/model/note"
	"neatly/internal/model/tag"
	"neatly/pkg/logging"
	"neatly/pkg/postgres"
)

type Authorisation interface {
	CreateAccount(u *auth.Account) error
	GetAccount(u *auth.Account) error
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
	Assign(tagID, noteID, userID int) error
	Update(userID, tagID int, t tag.Tag) error
}

type Repository struct {
	Authorisation
	Note
	Tag
}

func New(db *sqlx.DB, logger logging.Logger) *Repository {
	return &Repository{
		Authorisation: postgres.NewAuthPostgres(db, logger),
		Note:          postgres.NewNotePostgres(db, logger),
		Tag:           postgres.NewTagPostgres(db, logger),
	}
}
