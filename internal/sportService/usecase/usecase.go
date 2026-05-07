package usecase

import (
	"HSE/internal/sportService"
	"context"
	"errors"
	"strings"
)

type useCase struct {
	repo sportService.Repository // ← интерфейс из sportService
}

// New - конструктор
func New(repo sportService.Repository) sportService.UseCase { // ← возвращаем интерфейс UseCase
	return &useCase{repo: repo}
}

func (uc *useCase) GetTestMessage(ctx context.Context) (string, error) {
	return uc.repo.GetTestMessage(ctx)
}

func (uc *useCase) SaveMessage(ctx context.Context, message string) (string, error) {
	message = strings.TrimSpace(message)
	if message == "" {
		return "", errors.New("message is required")
	}

	return uc.repo.SaveMessage(ctx, message)
}

func (uc *useCase) ListUsers(ctx context.Context, limit int) ([]sportService.UserRecord, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	return uc.repo.ListUsers(ctx, limit)
}

func (uc *useCase) GetUserByID(ctx context.Context, id string) (*sportService.UserRecord, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, errors.New("id is required")
	}

	return uc.repo.GetUserByID(ctx, id)
}
