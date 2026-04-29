package errs

import "errors"

var (
	// common
	ErrInvalidInput = errors.New("invalid input")

	// not found
	ErrUserNotFound         = errors.New("user not found")
	ErrProductNotFound      = errors.New("product not found")
	ErrOrderNotFound        = errors.New("order not found")
	ErrOrderItemNotFound    = errors.New("order item not found")
	ErrProductSizeNotFound  = errors.New("product size not found")
	ErrProductImageNotFound = errors.New("product image not found")

	// postgres mapped
	ErrConflict     = errors.New("conflict")              // 23505
	ErrForeignKey   = errors.New("foreign key violation") // 23503
	ErrValueTooLong = errors.New("value too long")        // 22001
	ErrBadFormat    = errors.New("bad format")            // 22P02
	ErrNotNull      = errors.New("not null violation")    // 23502
)
