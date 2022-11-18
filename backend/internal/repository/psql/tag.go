package psql

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"neatly/internal/model"
	"neatly/pkg/client/psqlclient"
	"neatly/pkg/e"
	"neatly/pkg/logging"
)

type TagPostgres struct {
	db     *sqlx.DB
	logger logging.Logger
}

func NewTagPostgres(client *psqlclient.Client, logger logging.Logger) *TagPostgres {
	return &TagPostgres{db: client.DB, logger: logger}
}

func (r *TagPostgres) Create(userID, noteID int, t *model.Tag) error {

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	createTagQuery := `INSERT INTO tags AS t (label, color) VALUES ($1, $2) RETURNING id`

	r.logger.Infof("Tag with id %v created", t.ID)

	row := r.db.QueryRow(createTagQuery, t.Label, t.Color)
	err = row.Scan(&t.ID)

	if err != nil {
		return err
	}

	r.logger.Infof("Connecting tag with id %v and accounts with id with id %v", t.ID, userID)
	userTagQuery := `INSERT INTO users_tags (users_id, tags_id)
				    SELECT $1, $2 WHERE NOT EXISTS (
    			       SELECT users_id, tags_id FROM users_tags WHERE users_id = $1 AND tags_id = $2
    			)`
	_, err = tx.Exec(userTagQuery, userID, t.ID)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (r *TagPostgres) Assign(tagID, noteID, userID int) error {
	r.logger.Infof("Assigning tag with id %v to note with id with id %v", tagID, noteID)
	assignTagQuery := `INSERT INTO tags_notes (notes_id, tags_id) VALUES ($1, $2)`
	_, err := r.db.Exec(assignTagQuery, noteID, tagID)
	if err != nil {
		r.logger.Infof("Internal error: %v", err.Error())
		if err == sql.ErrNoRows {
			return e.ClientNoteError
		}
		return err
	}

	return nil
}

func (r *TagPostgres) GetAll(userID int) ([]model.Tag, error) {
	var tags []model.Tag
	tags = make([]model.Tag, 0)

	query := `SELECT tags_id AS id, label, color FROM
			  tags t JOIN users_tags ut ON ut.tags_id = t.id  WHERE
			  ut.users_id = $1`

	err := r.db.Select(&tags, query, userID)
	if err != nil {
		r.logger.Info(err)
	}
	return tags, err
}

func (r *TagPostgres) GetAllByNote(userID, noteID int) ([]model.Tag, error) {
	var tags []model.Tag
	tags = make([]model.Tag, 0)

	query := `SELECT t.id AS id, label, color FROM tags t
    		  JOIN users_tags ut ON ut.tags_id = t.id
    		  JOIN tags_notes nt on t.id = nt.tags_id
    		  WHERE users_id = $1 AND notes_id = $2`

	err := r.db.Select(&tags, query, userID, noteID)
	if err != nil {
		r.logger.Info(err)
	}
	return tags, err
}

func (r *TagPostgres) GetOne(userID, tagID int) (model.Tag, error) {
	var t model.Tag

	query := `SELECT t.id AS id, label, color FROM tags t
    		  JOIN users_tags ut ON ut.tags_id = t.id
    		  JOIN tags_notes nt on t.id = nt.tags_id
    		  WHERE users_id = $1 AND t.id = $2`

	err := r.db.Get(&t, query, userID, tagID)
	if err != nil {
		r.logger.Infof("Internal error: %v", err.Error())
		if err == sql.ErrNoRows {
			return t, e.ClientTagError
		}
	}
	return t, nil
}

func (r *TagPostgres) Delete(userID, tagID int) error {
	query := `DELETE FROM tags t USING users_tags ut WHERE 
              t.id = ut.tags_id AND ut.users_id = $1 AND ut.tags_id = $2`
	_, err := r.db.Exec(query, userID, tagID)

	return err
}

func (r *TagPostgres) Update(userID, tagID int, t model.Tag) error {
	query := `UPDATE tags t SET label=$1, color=$2 FROM users_tags ut
              WHERE t.id = ut.tags_id AND ut.tags_id = $3 AND ut.users_id = $4`
	_, err := r.db.Exec(query, t.Label, t.Color, tagID, userID)

	return err
}

func (r *TagPostgres) Detach(userID, tagID, noteID int) error {
	query := `DELETE FROM tags_notes USING users_tags ut WHERE
              tags_notes.tags_id = ut.tags_id AND ut.users_id = $1 AND ut.tags_id = $2 AND notes_id = $3`
	_, err := r.db.Exec(query, userID, tagID, noteID)

	return err
}
