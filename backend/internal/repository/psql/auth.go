package psql

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"neatly/internal/model/account"
	"neatly/pkg/client/psqlclient"
	"neatly/pkg/logging"
)

const (
	usersTable = "users"
)

type AuthPostgres struct {
	db     *sqlx.DB
	logger logging.Logger
}

func NewAuthPostgres(client *psqlclient.Client, logger logging.Logger) *AuthPostgres {
	return &AuthPostgres{
		db:     client.DB,
		logger: logger,
	}
}

func (r *AuthPostgres) CreateAccount(u *account.Account) error {
	query := fmt.Sprintf(
		"INSERT INTO %s (name, username, email, password_hash) VALUES ($1, $2, $3, $4) RETURNING id",
		usersTable,
	)

	r.logger.Info("Creating accounts")
	row := r.db.QueryRow(query, u.Name, u.Username, u.Email, u.PasswordHash)
	if err := row.Scan(&u.ID); err != nil {
		r.logger.Info(err)
		return &account.CanNotCreateAccountErr{}
	}
	return nil
}

func (r *AuthPostgres) AuthorizeAccount(u *account.Account) error {
	query := fmt.Sprintf(
		"SELECT id, name, username, password_hash, email FROM %s WHERE username=$1",
		usersTable,
	)

	r.logger.Infof("Get accounts %v from db", u.ID)
	err := r.db.Get(u, query, &u.Username)
	if err != nil {
		r.logger.Info(err)
		if errors.Is(err, sql.ErrNoRows) {
			return &account.AccountNotFoundErr{}
		}
		return &account.CanNotLoginErr{}
	}

	return nil
}

func (r *AuthPostgres) GetOne(userID int) (account.Account, error) {
	query := fmt.Sprintf(
		"SELECT id, name, username, email FROM %s WHERE id=$1",
		usersTable,
	)

	var a account.Account

	err := r.db.Get(&a, query, userID)
	if err != nil {
		r.logger.Info(err)
		if errors.Is(err, sql.ErrNoRows) {
			return a, &account.AccountNotFoundErr{}
		}
		return a, &account.CanNotGetErr{}
	}
	r.logger.Infof("Got account %v from db", a.ID)

	return a, nil
}
