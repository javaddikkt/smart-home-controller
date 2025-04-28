package postgres

import (
	"context"
	"errors"
	"fmt"
	"homework/internal/domain"
	"homework/internal/usecase"

	"github.com/jackc/pgx/v5"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool}
}

func (r *UserRepository) SaveUser(ctx context.Context, user *domain.User) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("user is nil")
	}

	const sql = `
        INSERT INTO users (name)
        VALUES ($1)
        RETURNING id
    `

	err := r.pool.QueryRow(ctx, sql, user.Name).Scan(&user.ID)
	return err
}

func (r *UserRepository) GetUserByID(ctx context.Context, id int64) (*domain.User, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	const sql = `
		SELECT id, name
		FROM users
		WHERE id = $1
  `

	row := r.pool.QueryRow(ctx, sql, id)

	var u domain.User
	if err := row.Scan(&u.ID, &u.Name); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, usecase.ErrUserNotFound
		}
		return nil, fmt.Errorf("user %d not found: %w", id, err)
	}
	return &u, nil
}
