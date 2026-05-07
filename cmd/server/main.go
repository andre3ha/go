package main

import (
	"context"
	"log/slog"
	stdhttp "net/http" // ← алиас для стандартного http
	"os"
	"os/signal"
	"syscall"
	"time"

	"HSE/internal/config"
	"HSE/internal/sportService/http"
	"HSE/internal/sportService/repository"
	"HSE/internal/sportService/usecase"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	if err := godotenv.Overload(); err != nil {
		logger.Error("failed to load .env", "error", err)
		os.Exit(1)
	}

	cfg, err := config.New()
	if err != nil {
		logger.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	db, err := sqlx.Connect("postgres", cfg.DB.DSN())
	if err != nil {
		logger.Error("failed to connect to db", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// Создание зависимостей
	repo := repository.New(db)
	uc := usecase.New(repo)
	h := http.New(uc) // ← h реализует интерфейс Handler

	// Регистрация маршрута
	mux := stdhttp.NewServeMux() // ← стандартный http через алиас
	mux.HandleFunc("/test", h.Test())
	mux.HandleFunc("/dbtest", h.DBTest())
	mux.HandleFunc("/users", h.Users())
	mux.HandleFunc("/users/", h.UserByID())

	// Настройка сервера
	srv := &stdhttp.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: withCORS(mux),
	}

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logger.Info("server started", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != stdhttp.ErrServerClosed {
			logger.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	<-quit
	logger.Info("shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("server shutdown error", "error", err)
		os.Exit(1)
	}

	logger.Info("server stopped gracefully")
}

func withCORS(next stdhttp.Handler) stdhttp.Handler {
	return stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == stdhttp.MethodOptions {
			w.WriteHeader(stdhttp.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
