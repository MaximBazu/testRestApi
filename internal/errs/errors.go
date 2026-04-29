package errs

import "errors"

var (
	ErrUserNotFound         = errors.New("user not found")
	ErrProductNotFound      = errors.New("product not found")
	ErrOrderNotFound        = errors.New("order not found")
	ErrOrderItemNotFound    = errors.New("order item not found")
	ErrProductSizeNotFound  = errors.New("product size not found")
	ErrProductImageNotFound = errors.New("product image not found")
	ErrInvalidInput         = errors.New("invalid input")

	// DB-level mapped errors
	ErrConflict     = errors.New("conflict")
	ErrForeignKey   = errors.New("foreign key violation")
	ErrValueTooLong = errors.New("value too long")
	ErrBadFormat    = errors.New("bad format")
	ErrNotNull      = errors.New("not null violation")
)
