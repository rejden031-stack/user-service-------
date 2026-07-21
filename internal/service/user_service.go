package service

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type UserService struct {
	pool *pgxpool.Pool
}

func NewUserService(pool *pgxpool.Pool) *UserService {
	return &UserService{pool: pool}
}

func (s *UserService) GetUser(ctx context.Context, id int) (User, error) {
	var user User
	err := s.pool.QueryRow(ctx, "SELECT id, name FROM users WHERE id = $1", id).Scan(&user.ID, &user.Name)
	if err != nil {
		return User{}, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

func (s *UserService) GetAllUsers(ctx context.Context) ([]User, error) {
	rows, err := s.pool.Query(ctx, "SELECT id, name FROM users")
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}
	return users, nil
}

func (s *UserService) CreateUser(ctx context.Context, name string) (User, error) {
	var user User
	err := s.pool.QueryRow(ctx, "INSERT INTO users (name) VALUES ($1) RETURNING id, name", name).Scan(&user.ID, &user.Name)
	if err != nil {
		return User{}, fmt.Errorf("failed to create user: %w", err)
	}
	return user, nil
}
