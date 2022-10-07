package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"neatly/internal/model/tag"
	"neatly/pkg/e"
	"neatly/pkg/logging"
)

const (
	tagsTable      = "tags"
	notesTagsTable = "tags_notes"
	usersTagsTable = "users_tags"
)

type TagPostgres struct {
	db     *sqlx.DB
	logger logging.Logger
}

func NewTagPostgres(db *sqlx.DB, logger logging.Logger) *TagPostgres {
	return &TagPostgres{db: db, logger: logger}
}

func (r *TagPostgres) Create(userID, noteID int, t *tag.Tag) error {

	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	createTagQuery := fmt.Sprintf(
		`INSERT INTO %s AS t (name, color) VALUES ($1, $2) RETURNING id`,
		tagsTable)

	r.logger.Infof("Tag with id %v created", t.ID)

	row := r.db.QueryRow(createTagQuery, t.Name, t.Color)
	err = row.Scan(&t.ID)

	if err != nil {
		r.logger.Info(err)
		if errors.Is(err, sql.ErrNoRows) {
			return &e.CanNotCreateNoteErr{}
		}
		return err
	}

	r.logger.Infof("Connecting tag with id %v and auth with id with id %v", t.ID, userID)
	userTagQuery := fmt.Sprintf(
		`INSERT INTO %s
    			(users_id, tags_id)
				SELECT $1, $2
				WHERE
    			NOT EXISTS (
    			    SELECT users_id, tags_id FROM users_tags WHERE users_id = $1 AND tags_id = $2
    			);`, usersTagsTable)
	_, err = tx.Exec(userTagQuery, userID, t.ID)
	if err != nil {
		tx.Rollback()
		r.logger.Info(err)
		if errors.Is(err, sql.ErrNoRows) {
			return &e.CanNotCreateNoteErr{}
		}
		return err
	}
	return tx.Commit()
}

func (r *TagPostgres) Assign(tagID, noteID, userID int) error {
	r.logger.Infof("Assigning tag with id %v to note with id with id %v", tagID, noteID)
	assignTagQuery := fmt.Sprintf(
		`INSERT INTO %s (notes_id, tags_id) VALUES ($1, $2)`, notesTagsTable)
	_, err := r.db.Exec(assignTagQuery, noteID, tagID)
	if err != nil {
		r.logger.Info(err)
		return err
	}

	return nil
}

func (r *TagPostgres) GetAll(userID int) ([]tag.Tag, error) {
	var tags []tag.Tag

	query := fmt.Sprintf(`SELECT tags_id AS id, name, color FROM
								%s t INNER JOIN %s ut ON ut.tags_id = t.id  WHERE
								ut.users_id = $1`, tagsTable, usersTagsTable)

	err := r.db.Select(&tags, query, userID)
	if err != nil {
		r.logger.Info(err)
	}
	return tags, err
}

func (r *TagPostgres) GetAllByNote(userID, noteID int) ([]tag.Tag, error) {
	var tags []tag.Tag

	query := fmt.Sprintf(`SELECT t.id AS id, name, color FROM %s t
    							INNER JOIN %s ut ON ut.tags_id = t.id
    							INNER JOIN %s nt on t.id = nt.tags_id
    							WHERE users_id = $1 AND notes_id = $2`,
		tagsTable, usersTagsTable, notesTagsTable)

	err := r.db.Select(&tags, query, userID, noteID)
	if err != nil {
		r.logger.Info(err)
	}
	return tags, err
}

func (r *TagPostgres) GetOne(userID, tagID int) (tag.Tag, error) {
	var tag tag.Tag

	query := fmt.Sprintf(`SELECT t.id AS id, name, color FROM %s t
    							INNER JOIN %s ut ON ut.tags_id = t.id
    							INNER JOIN %s nt on t.id = nt.tags_id
    							WHERE users_id = $1 AND t.id = $2`,
		tagsTable, usersTagsTable, notesTagsTable)

	err := r.db.Get(&tag, query, userID, tagID)
	if err != nil {
		r.logger.Info(err)
		if errors.Is(err, sql.ErrNoRows) {
			return tag, &e.TagNotFoundErr{}
		}
	}
	return tag, err
}

func (r *TagPostgres) Delete(userID, tagID int) error {
	query := fmt.Sprintf(
		`DELETE FROM %s t USING %s ut WHERE 
              t.id = ut.tags_id AND ut.users_id = $1 AND ut.tags_id = $2`,
		tagsTable, usersTagsTable)
	_, err := r.db.Exec(query, userID, tagID)

	return err
}

func (r *TagPostgres) Update(userID, tagID int, t tag.Tag) error {
	query := fmt.Sprintf(
		`UPDATE %s t SET 
                name=$1, color=$2 FROM
                %s ut WHERE t.id = ut.tags_id AND 
				ut.tags_id = $3 AND ut.users_id = $4`,
		tagsTable, usersTagsTable)
	_, err := r.db.Exec(query, t.Name, t.Color, tagID, userID)

	return err
}
