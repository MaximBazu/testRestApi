package postgres

import (
	"errors"

	"RESTAPI/internal/errs"

	"github.com/jackc/pgx/v5/pgconn"
)

// MapPGError converts PostgreSQL/pgx errors to domain errors.
func MapPGError(err error) error {
	if err == nil {
		return nil
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		// unique_violation
		case "23505":
			return errs.ErrConflict
		// foreign_key_violation
		case "23503":
			return errs.ErrForeignKey
		// string_data_right_truncation
		case "22001":
			return errs.ErrValueTooLong
		// invalid_text_representation (e.g. bad UUID/int format)
		case "22P02":
			return errs.ErrBadFormat
		// not_null_violation
		case "23502":
			return errs.ErrNotNull
		}
	}

	return err
}
