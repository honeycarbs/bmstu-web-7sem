package postgres

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"neatly/internal/model/note"
	"neatly/pkg/logging"
)

const (
	notesTable      = "note"
	notesBodyTable  = "note_body"
	usersNotesTable = "users_notes"
)

type NotePostgres struct {
	db     *sqlx.DB
	logger logging.Logger
}

func NewNotePostgres(db *sqlx.DB, logger logging.Logger) *NotePostgres {
	return &NotePostgres{
		db:     db,
		logger: logger,
	}
}

func (r *NotePostgres) Create(userID int, n *note.Note) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	createNoteQuery := fmt.Sprintf(`
	INSERT INTO %s (header, short_body, color)
	VALUES ($1, $2, $3) RETURNING id`, notesTable)
	row := tx.QueryRow(createNoteQuery, n.Header, n.ShortBody, n.Color)
	if err := row.Scan(&n.ID); err != nil {
		tx.Rollback()
		r.logger.Error(err)
		return err
	}
	createNoteBodyQuery := fmt.Sprintf("INSERT INTO %s (id, body) VALUES ($1, $2)", notesBodyTable)
	_, err = tx.Exec(createNoteBodyQuery, n.ID, n.Body)
	if err != nil {
		tx.Rollback()
		r.logger.Error(err)
		return err
	}
	createUsersNoteQuery := fmt.Sprintf("INSERT INTO %s (user_id, note_id) VALUES ($1, $2)", usersNotesTable)
	_, err = tx.Exec(createUsersNoteQuery, userID, n.ID)
	if err != nil {
		tx.Rollback()
		r.logger.Error(err)
		return err
	}

	return tx.Commit()
}

func (r *NotePostgres) GetAll(userID int) ([]note.Note, error) {
	var notes []note.Note

	getNotesQuery := fmt.Sprintf(
		`SELECT n.id, n.header, n.short_body, n.color FROM %s n
    			INNER JOIN %s un ON n.id = un.note_id
    			WHERE un.user_id = $1`,
		notesTable,
		usersNotesTable,
	)

	err := r.db.Select(&notes, getNotesQuery, userID)

	if err != nil {
		r.logger.Info(err)
		return notes, err
	}

	return notes, err
}

func (r *NotePostgres) GetOne(userID, noteID int) (note.Note, error) {
	r.logger.Infof("USER ID: %v, NOTE ID: %v", userID, noteID)
	tx, err := r.db.Begin()
	if err != nil {
		return note.Note{}, err
	}
	var n note.Note

	selectNoteQuery := fmt.Sprintf(
		`SELECT n.id, n.header, n.short_body, n.color FROM
				%s n INNER JOIN %s un ON n.id = un.note_id
				WHERE un.user_id = $1 AND un.note_id = $2`,
		notesTable,
		usersNotesTable,
	)
	err = r.db.Get(&n, selectNoteQuery, userID, noteID)
	if err != nil {
		tx.Rollback()
		r.logger.Info(err)
		return note.Note{}, err
	}

	selectBodyQuery := fmt.Sprintf(
		`SELECT nb.body FROM %s nb INNER JOIN %s n ON nb.id = n.id
				WHERE n.id = $1`,
		notesBodyTable,
		notesTable,
	)

	err = r.db.Get(&n.Body, selectBodyQuery, noteID)
	if err != nil {
		tx.Rollback()
		r.logger.Info(err)
		return note.Note{}, err
	}

	return n, tx.Commit()
}

func (r *NotePostgres) Delete(userID, noteID int) error {
	query := fmt.Sprintf(
		`DELETE FROM %s n USING %s un WHERE 
              n.id = un.note_id AND un.user_id = $1 AND un.note_id = $2`,
		notesTable, usersNotesTable)
	_, err := r.db.Exec(query, userID, noteID)

	return err
}

func (r *NotePostgres) Update(userID int, n note.Note) error {
	// TODO: update tasks
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	noteQuery := fmt.Sprintf(
		`UPDATE %s n SET 
                header=$1, short_body=$2, color = $3 FROM
                %s un WHERE n.id = un.note_id AND 
				un.note_id = $4 AND un.user_id = $5`,
		notesTable, usersNotesTable)
	_, err = r.db.Exec(noteQuery, n.Header, n.ShortBody, n.Color, n.ID, userID)
	if err != nil {
		tx.Rollback()
		r.logger.Info(err)
		return err
	}

	bodyQuery := fmt.Sprintf(
		`UPDATE %s nb SET body=$2 WHERE nb.id = $1`,
		notesBodyTable)
	_, err = r.db.Exec(bodyQuery, n.ID, n.Body)
	if err != nil {
		tx.Rollback()
		r.logger.Info(err)
		return err
	}

	return tx.Commit()
}
