package sportService

import (
	"context"
	"net/http"
)

// Handler — контракт HTTP-хэндлеров
type Handler interface {
	Test() http.HandlerFunc
}

// UseCase — контракт бизнес-логики
type UseCase interface {
	GetTestMessage(ctx context.Context) (string, error)
}

// Repository — контракт работы с данными
type Repository interface {
	GetTestMessage(ctx context.Context) (string, error)
}
