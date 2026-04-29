package postgres

import (
	"RESTAPI/internal/errs"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

func MapPGError(err error) error {
	if err == nil {
		return nil
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23505":
			return errs.ErrConflict
		case "23503":
			return errs.ErrForeignKey
		case "22001":
			return errs.ErrValueTooLong
		case "22P02":
			return errs.ErrBadFormat
		case "23502":
			return errs.ErrNotNull
		}
	}

	return err
}
