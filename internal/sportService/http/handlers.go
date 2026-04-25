package http

import (
	"HSE/internal/sportService" // ← для интерфейса Handler
	"net/http"
)

type handler struct {
	uc sportService.UseCase // ← используем интерфейс из sportService
}

// New - конструктор
func New(uc sportService.UseCase) sportService.Handler { // ← возвращаем интерфейс Handler
	return &handler{uc: uc}
}

// Test - реализует метод интерфейса Handler
func (h *handler) Test() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		msg, err := h.uc.GetTestMessage(r.Context())
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(msg))
	}
}
