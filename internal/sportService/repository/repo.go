package repository

import (
	"HSE/internal/sportService"
	"context"
)

type repo struct{}

// New - конструктор
func New() sportService.Repository { // ← возвращаем интерфейс Repository
	return &repo{}
}

func (r *repo) GetTestMessage(ctx context.Context) (string, error) {
	return "Hello", nil
}
