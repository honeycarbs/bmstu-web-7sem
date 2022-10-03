package postgres

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"neatly/internal/model/tag"
	"neatly/pkg/logging"
)

const (
	tagsTable      = "tags"
	notesTagsTable = "notes_tags"
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
		return err
	}

	r.logger.Infof("Connecting tag with id %v and auth with id with id %v", t.ID, userID)
	userTagQuery := fmt.Sprintf(
		`INSERT INTO %s
    			(user_id, tag_id)
				SELECT $1, $2
				WHERE
    			NOT EXISTS (
    			    SELECT user_id, tag_id FROM users_tags WHERE user_id = $1 AND tag_id = $2
    			);`, usersTagsTable)
	_, err = tx.Exec(userTagQuery, userID, t.ID)
	if err != nil {
		tx.Rollback()
		r.logger.Info(err)
		return err
	}
	return tx.Commit()
}

func (r *TagPostgres) Assign(tagID, noteID, userID int) error {
	r.logger.Infof("Assigning tag with id %v to note with id with id %v", tagID, noteID)
	assignTagQuery := fmt.Sprintf(
		`INSERT INTO %s (note_id, tag_id) VALUES ($1, $2)`, notesTagsTable)
	_, err := r.db.Exec(assignTagQuery, noteID, tagID)
	if err != nil {
		r.logger.Info(err)
		return err
	}

	return nil
}

func (r *TagPostgres) GetAll(userID int) ([]tag.Tag, error) {
	var tags []tag.Tag

	query := fmt.Sprintf(`SELECT tag_id AS id, name, color FROM
								%s t INNER JOIN %s ut ON ut.tag_id = t.id  WHERE
								ut.user_id = $1`, tagsTable, usersTagsTable)

	err := r.db.Select(&tags, query, userID)
	if err != nil {
		r.logger.Info(err)
	}
	return tags, err
}

func (r *TagPostgres) GetAllByNote(userID, noteID int) ([]tag.Tag, error) {
	var tags []tag.Tag

	query := fmt.Sprintf(`SELECT t.id AS id, name, color FROM %s t
    							INNER JOIN %s ut ON ut.tag_id = t.id
    							INNER JOIN %s nt on t.id = nt.tag_id
    							WHERE user_id = $1 AND note_id = $2`,
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
    							INNER JOIN %s ut ON ut.tag_id = t.id
    							INNER JOIN %s nt on t.id = nt.tag_id
    							WHERE user_id = $1 AND t.id = $2`,
		tagsTable, usersTagsTable, notesTagsTable)

	err := r.db.Get(&tag, query, userID, tagID)
	if err != nil {
		r.logger.Info(err)
	}
	return tag, err
}

func (r *TagPostgres) Delete(userID, tagID int) error {
	query := fmt.Sprintf(
		`DELETE FROM %s t USING %s ut WHERE 
              t.id = ut.tag_id AND ut.user_id = $1 AND ut.tag_id = $2`,
		tagsTable, usersTagsTable)
	_, err := r.db.Exec(query, userID, tagID)

	return err
}

func (r *TagPostgres) Update(userID, tagID int, t tag.Tag) error {
	query := fmt.Sprintf(
		`UPDATE %s t SET 
                name=$1, color=$2 FROM
                %s ut WHERE t.id = ut.tag_id AND 
				ut.tag_id = $3 AND ut.user_id = $4`,
		tagsTable, usersTagsTable)
	_, err := r.db.Exec(query, t.Name, t.Color, tagID, userID)

	return err
}
