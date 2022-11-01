package psql

import (
	"github.com/jmoiron/sqlx"
	"neatly/internal/model"
	"neatly/pkg/client/psqlclient"
	"neatly/pkg/logging"
)

type AuthPostgres struct {
	db     *sqlx.DB
	logger logging.Logger
}

func NewAccountPostgres(client *psqlclient.Client, logger logging.Logger) *AuthPostgres {
	return &AuthPostgres{
		db:     client.DB,
		logger: logger,
	}
}

func (r *AuthPostgres) CreateAccount(a *model.Account) error {
	query := `INSERT INTO users
              (name, username, email, password_hash)
              VALUES ($1, $2, $3, $4) RETURNING id`

	r.logger.Info("Creating accounts")
	row := r.db.QueryRow(query, a.Name, a.Username, a.Email, a.PasswordHash)
	if err := row.Scan(&a.ID); err != nil {
		return err
	}
	return nil
}

func (r *AuthPostgres) AuthorizeAccount(a *model.Account) error {
	query := `SELECT id, name, username, password_hash, email
			  FROM users WHERE username=$1`

	err := r.db.Get(a, query, &a.Username)
	if err != nil {
		return err
	}

	return nil
}

func (r *AuthPostgres) GetOne(userID int) (model.Account, error) {
	query := `SELECT id, name, username, email FROM
              users WHERE id=$1`

	var a model.Account

	err := r.db.Get(&a, query, userID)
	if err != nil {
		return a, err
	}
	return a, nil
}
