package postgres

import (
	"RESTAPI/internal/errs"
	"context"
	"errors"
	"fmt"
	"time"

	"RESTAPI/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DBTX interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
}
type userRepository struct {
	db DBTX
}

func NewUserRepository(db *pgxpool.Pool) *userRepository {
	return &userRepository{db: db}
}

func NewUserRepositoryTx(tx pgx.Tx) *userRepository {
	return &userRepository{db: tx}
}

func (r *userRepository) GetByID(ctx context.Context, id int) (*model.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	query := `SELECT id, name, surname, email, telegram_tag, created_at FROM users WHERE id=$1`

	var user model.User

	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Surname,
		&user.Email,
		&user.TelegramTag,
		&user.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("%w: %v", errs.ErrUserNotFound, err)
		}
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	query := `
		INSERT INTO users (name, surname, email, telegram_tag)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`

	return r.db.QueryRow(ctx, query,
		ctx,
		query,
		user.Name,
		user.Surname,
		user.Email,
		user.TelegramTag,
	).Scan(&user.ID, &user.CreatedAt)
}
