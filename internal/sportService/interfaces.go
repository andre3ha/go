package sportService

import (
	"context"
	"net/http"
	"time"
)

type UserRecord struct {
	ID        string    `db:"id" json:"id"`
	Login     string    `db:"login" json:"login"`
	Email     string    `db:"email" json:"email"`
	Phone     string    `db:"phone" json:"phone"`
	IsActive  bool      `db:"is_active" json:"is_active"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

// Repository describes data access methods for service.
type Repository interface {
	GetTestMessage(ctx context.Context) (string, error)
	SaveMessage(ctx context.Context, message string) (string, error)
	ListUsers(ctx context.Context, limit int) ([]UserRecord, error)
	GetUserByID(ctx context.Context, id string) (*UserRecord, error)
}

// UseCase describes business logic methods.
type UseCase interface {
	GetTestMessage(ctx context.Context) (string, error)
	SaveMessage(ctx context.Context, message string) (string, error)
	ListUsers(ctx context.Context, limit int) ([]UserRecord, error)
	GetUserByID(ctx context.Context, id string) (*UserRecord, error)
}

// Handler describes HTTP handlers provided by service.
type Handler interface {
	Test() http.HandlerFunc
	DBTest() http.HandlerFunc
	Users() http.HandlerFunc
	UserByID() http.HandlerFunc
}
