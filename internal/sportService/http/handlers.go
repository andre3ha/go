package http

import (
	"HSE/internal/sportService" // ← для интерфейса Handler
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
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

func (h *handler) DBTest() http.HandlerFunc {
	type request struct {
		Message string `json:"message"`
	}
	type response struct {
		Status string `json:"status"`
		ID     string `json:"id"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			var req request
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "invalid json body", http.StatusBadRequest)
				return
			}

			id, err := h.uc.SaveMessage(r.Context(), req.Message)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(response{
				Status: "saved",
				ID:     id,
			})
		case http.MethodGet:
			users, err := h.uc.ListUsers(r.Context(), 10)
			if err != nil {
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(users)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
	}
}

func (h *handler) Users() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		limit := 10
		if rawLimit := r.URL.Query().Get("limit"); rawLimit != "" {
			parsed, err := strconv.Atoi(rawLimit)
			if err != nil {
				http.Error(w, "invalid limit", http.StatusBadRequest)
				return
			}
			limit = parsed
		}

		users, err := h.uc.ListUsers(r.Context(), limit)
		if err != nil {
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(users)
	}
}

func (h *handler) UserByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		id := strings.TrimPrefix(r.URL.Path, "/users/")
		if id == "" || strings.Contains(id, "/") {
			http.Error(w, "invalid user id", http.StatusBadRequest)
			return
		}

		user, err := h.uc.GetUserByID(r.Context(), id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if user == nil {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(user)
	}
}
