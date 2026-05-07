package repository

import (
	"HSE/internal/sportService"
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type repo struct {
	db *sqlx.DB
}

// New - конструктор
func New(db *sqlx.DB) sportService.Repository { // ← возвращаем интерфейс Repository
	return &repo{db: db}
}

func (r *repo) GetTestMessage(ctx context.Context) (string, error) {
	const query = "SELECT COUNT(*) FROM users"

	var usersCount int
	if err := r.db.GetContext(ctx, &usersCount, query); err != nil {
		return "", err
	}

	return fmt.Sprintf("db is connected, users_count=%d", usersCount), nil
}

func (r *repo) SaveMessage(ctx context.Context, message string) (string, error) {
	login := fmt.Sprintf("dbtest_%d", time.Now().UnixNano())
	email := fmt.Sprintf("%s@example.local", login)
	password := "dbtest_password"

	const query = `
		INSERT INTO users (login, password, email, phone, is_active)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	var userID string
	if err := r.db.GetContext(ctx, &userID, query, login, password, email, message, true); err != nil {
		return "", err
	}

	return userID, nil
}

func (r *repo) ListUsers(ctx context.Context, limit int) ([]sportService.UserRecord, error) {
	const query = `
		SELECT id, login, email, phone, is_active, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
		LIMIT $1
	`

	var users []sportService.UserRecord
	if err := r.db.SelectContext(ctx, &users, query, limit); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *repo) GetUserByID(ctx context.Context, id string) (*sportService.UserRecord, error) {
	const query = `
		SELECT id, login, email, phone, is_active, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var user sportService.UserRecord
	if err := r.db.GetContext(ctx, &user, query, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}
