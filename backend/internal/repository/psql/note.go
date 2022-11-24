package psql

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"neatly/internal/model"
	"neatly/pkg/dbclient"
	"neatly/pkg/e"
	"neatly/pkg/logging"
	"time"
)

type NotePostgres struct {
	db     *sqlx.DB
	logger logging.Logger
}

func NewNotePostgres(client *dbclient.Client, logger logging.Logger) *NotePostgres {
	return &NotePostgres{
		db:     client.DB,
		logger: logger,
	}
}

func (r *NotePostgres) Create(userID int, n *model.Note) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	createNoteQuery := `INSERT INTO notes (header, short_body, color, edited)
						VALUES ($1, $2, $3, $4) RETURNING id`

	row := tx.QueryRow(createNoteQuery, n.Header, n.ShortBody, n.Color, time.Now())
	if err := row.Scan(&n.ID); err != nil {
		tx.Rollback()
		r.logger.Error(err)
		return e.InternalDBError
	}
	createNoteBodyQuery := `INSERT INTO notes_body (id, body) VALUES ($1, $2)`
	_, err = tx.Exec(createNoteBodyQuery, n.ID, n.Body)
	if err != nil {
		tx.Rollback()
		r.logger.Error(err)
		return e.InternalDBError
	}
	createUsersNoteQuery := `INSERT INTO users_notes (users_id, notes_id) VALUES ($1, $2)`
	_, err = tx.Exec(createUsersNoteQuery, userID, n.ID)
	if err != nil {
		tx.Rollback()
		r.logger.Error(err)
		return e.InternalDBError
	}

	return tx.Commit()
}

func (r *NotePostgres) GetAll(userID int) ([]model.Note, error) {
	var notes []model.Note
	notes = make([]model.Note, 0)

	getNotesQuery := `SELECT n.id, n.header, n.short_body, n.color, n.edited FROM notes n
    			      JOIN users_notes un ON n.id = un.notes_id
    			      WHERE un.users_id = $1`

	err := r.db.Select(&notes, getNotesQuery, userID)
	if err != nil {
		r.logger.Info(err)
		return notes, err
	}

	return notes, err
}

func (r *NotePostgres) GetOne(userID, noteID int) (model.Note, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return model.Note{}, err
	}
	var n model.Note

	selectNoteQuery := `SELECT n.id, n.header, n.short_body, n.color, n.edited FROM
				        notes n JOIN users_notes un ON n.id = un.notes_id
				        WHERE un.users_id = $1 AND un.notes_id = $2`

	err = r.db.Get(&n, selectNoteQuery, userID, noteID)
	if err != nil {
		tx.Rollback()
		r.logger.Infof("Internal error: %v", err.Error())
		if err == sql.ErrNoRows {
			return n, e.ClientNoteError
		}
	}

	selectBodyQuery := `SELECT nb.body FROM notes_body nb JOIN notes n ON nb.id = n.id
				        WHERE n.id = $1`

	err = r.db.Get(&n.Body, selectBodyQuery, noteID)
	if err != nil {
		tx.Rollback()
		return n, err
	}

	return n, tx.Commit()
}

func (r *NotePostgres) Delete(userID, noteID int) error {
	query := `DELETE FROM notes WHERE 
              notes.id = $2`
	_, err := r.db.Exec(query, userID, noteID)

	return err
}

func (r *NotePostgres) Update(userID int, n model.Note) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	noteQuery := `UPDATE notes SET 
                  header=$1, short_body=$2, color = $3, edited=$4 FROM
                  users_notes WHERE notes.id = users_notes.notes_id AND 
				  users_notes.notes_id = $5 AND users_notes.users_id = $6`
	_, err = r.db.Exec(
		noteQuery,
		n.Header,
		n.ShortBody,
		n.Color,
		time.Now().UTC().Format(time.RFC3339),
		n.ID,
		userID,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	bodyQuery := `UPDATE notes_body SET body=$2 WHERE notes_body.id = $1`
	_, err = r.db.Exec(bodyQuery, n.ID, n.Body)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
