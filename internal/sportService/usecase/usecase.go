package usecase

import (
	"HSE/internal/sportService"
	"context"
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
