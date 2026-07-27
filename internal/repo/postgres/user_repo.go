package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"user-service/internal/entity"
)

type UserRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{pool: pool}
}

func (r *UserRepo) GetByID(ctx context.Context, id int) (entity.User, error) {
	var u entity.User
	err := r.pool.QueryRow(ctx, "SELECT id, name FROM users WHERE id = $1", id).Scan(&u.ID, &u.Name)
	if err != nil {
		return entity.User{}, fmt.Errorf("failed to get user: %w", err)
	}
	return u, nil
}

func (r *UserRepo) GetAll(ctx context.Context) ([]entity.User, error) {
	rows, err := r.pool.Query(ctx, "SELECT id, name FROM users")
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	defer rows.Close()

	var users []entity.User
	for rows.Next() {
		var u entity.User
		if err := rows.Scan(&u.ID, &u.Name); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, u)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return users, nil
}

func (r *UserRepo) Create(ctx context.Context, name string) (entity.User, error) {
	var u entity.User
	err := r.pool.QueryRow(ctx, "INSERT INTO users (name) VALUES ($1) RETURNING id, name", name).Scan(&u.ID, &u.Name)
	if err != nil {
		return entity.User{}, fmt.Errorf("failed to create user: %w", err)
	}
	return u, nil
}
