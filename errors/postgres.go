package errors

import (
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
)

const (
	duplicateKeyErrorCode = "23505"
)

func IsPostgresDuplicateKeyError(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == duplicateKeyErrorCode {
		return true
	}

	return false
}
