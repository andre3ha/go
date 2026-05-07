package domain

import "time"

// ========== Модели данных ==========

type User struct {
	ID        int       `json:"id"`
	Login     string    `json:"login"`
	Password  string    `json:"-"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Session struct {
	SessionID string    `json:"session_id"`
	UserID    int       `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

type Order struct {
	ID          string     `json:"id"`
	UserID      int        `json:"user_id"`
	Status      string     `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type OrderItem struct {
	OrderID   string  `json:"order_id"`
	ProductID int     `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

// ========== Repository интерфейсы ==========

type UserRepository interface {
	CreateUser(user *User) error
	GetUserByID(id int) (*User, error)
	GetUserByEmail(email string) (*User, error) // для логина
	GetUserByLogin(login string) (*User, error) // для регистрации
}

type SessionRepository interface {
	CreateSession(session *Session) error
	GetSessionByToken(token string) (*Session, error)
	DeleteSession(token string) error
}

type OrderRepository interface {
	CreateOrder(order *Order) error
	GetOrdersByUserID(userID int) ([]Order, error)
	GetOrderByID(id string) (*Order, error)
	UpdateOrderStatus(id string, status string) error
}

// ========== Usecase интерфейсы ==========

type UserUsecase interface {
	Register(login, password, email string) (*User, error)
	Login(login, password string) (*Session, error)
	Logout(token string) error
}

type OrderUsecase interface {
	PlaceOrder(userID int, items []OrderItem) (*Order, error)
	GetOrders(userID int) ([]Order, error)
}
