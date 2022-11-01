package psql

import (
	"github.com/lib/pq"
	"neatly/pkg/e"
	"neatly/pkg/logging"
)

const (
	codeDuplicateVal   = "23505"
	usernameConstraint = "users_username_key"
)

func ParsePsqlError(err *pq.Error) error {
	logging.GetLogger().Info(err.Code, err.Constraint)
	switch {
	case err.Code == codeDuplicateVal && err.Constraint == usernameConstraint:
		return e.ClientAccountError
	default:
		return err
	}
}
