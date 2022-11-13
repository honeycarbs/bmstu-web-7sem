package psql

import (
	"github.com/lib/pq"
	"neatly/pkg/e"
)

const (
	codeDuplicateVal   = "23505"
	usernameConstraint = "users_username_key"
)

func ParsePsqlError(err error) error {
	pqerr, ok := err.(*pq.Error)
	if !ok {
		return err
	}
	switch {
	case pqerr.Code == codeDuplicateVal && pqerr.Constraint == usernameConstraint:
		return e.ClientAccountError
	default:
		return err
	}
}
