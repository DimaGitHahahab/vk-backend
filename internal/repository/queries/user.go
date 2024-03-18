package queries

import (
	"context"
	"fmt"
	"vk-backend/internal/domain"
)

const addUser = `INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id, username, password, role`

func (q *Queries) AddUser(ctx context.Context, name string, password string) (*domain.User, error) {
	row := q.pool.QueryRow(ctx, addUser, name, password)
	user := &domain.User{}
	if err := row.Scan(&user.Id, &user.Name, &user.Password, &user.Role); err != nil {
		return nil, fmt.Errorf("failed to add user: %w", err)
	}
	return user, nil
}

const userExists = `SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)`

func (q *Queries) UserExists(ctx context.Context, name string) (bool, error) {
	var exists bool
	if err := q.pool.QueryRow(ctx, userExists, name).Scan(&exists); err != nil {
		return false, fmt.Errorf("failed to check if user exists: %w", err)
	}
	return exists, nil
}

const getUserByName = `SELECT id, username, password, role FROM users WHERE username = $1`

func (q *Queries) GetUserByName(ctx context.Context, name string) (*domain.User, error) {
	row := q.pool.QueryRow(ctx, getUserByName, name)
	user := &domain.User{}
	if err := row.Scan(&user.Id, &user.Name, &user.Password, &user.Role); err != nil {
		return nil, fmt.Errorf("failed to get user by name: %w", err)
	}
	return user, nil
}

const getUserById = `SELECT id, username, password, role FROM users WHERE id = $1`

func (q *Queries) GetUserById(ctx context.Context, id int) (*domain.User, error) {
	row := q.pool.QueryRow(ctx, getUserById, id)
	user := &domain.User{}
	if err := row.Scan(&user.Id, &user.Name, &user.Password, &user.Role); err != nil {
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}
	return user, nil
}
