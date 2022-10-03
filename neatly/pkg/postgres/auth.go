package postgres

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"neatly/internal/model/user"
	"neatly/pkg/logging"
)

const (
	usersTable = "users"
)

type AuthPostgres struct {
	db     *sqlx.DB
	logger logging.Logger
}

func NewAuthPostgres(db *sqlx.DB, logger logging.Logger) *AuthPostgres {
	return &AuthPostgres{
		db:     db,
		logger: logger,
	}
}

func (r *AuthPostgres) CreateUser(u *user.User) error {
	query := fmt.Sprintf(
		"INSERT INTO %s (name, username, email, password_hash) VALUES ($1, $2, $3, $4) RETURNING id",
		usersTable,
	)

	r.logger.Info("Creating user")
	row := r.db.QueryRow(query, u.Name, u.Username, u.Email, u.PasswordHash)
	if err := row.Scan(&u.ID); err != nil {
		r.logger.Info(err)
		return err
	}
	return nil
}

func (r *AuthPostgres) GetUser(u *user.User) error {
	logging.GetLogger().Infof("%s, %s", u.Username, u.Password)
	query := fmt.Sprintf(
		"SELECT id, password_hash FROM %s WHERE username=$1",
		usersTable,
	)

	r.logger.Infof("Get user %v from db", u.ID)
	err := r.db.Get(u, query, &u.Username)
	if err != nil {
		r.logger.Info(err)
		return err
	}

	r.logger.Infof("Сhecking users %v password", u.ID)
	err = u.CheckPassword(u.Password)
	if err != nil {
		r.logger.Info(err)
		return err
	}

	return nil
}
