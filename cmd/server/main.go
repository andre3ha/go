package main

import (
	"context"
	"log"
	stdhttp "net/http" // ← алиас для стандартного http
	"os"
	"os/signal"
	"syscall"
	"time"

	"HSE/internal/sportService/http"
	"HSE/internal/sportService/repository"
	"HSE/internal/sportService/usecase"
)

func main() {
	// Создание зависимостей
	repo := repository.New()
	uc := usecase.New(repo)
	h := http.New(uc) // ← h реализует интерфейс Handler

	// Регистрация маршрута
	mux := stdhttp.NewServeMux() // ← стандартный http через алиас
	mux.HandleFunc("/test", h.Test())

	// Настройка сервера
	srv := &stdhttp.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Println("Server started on http://localhost:8080/")
		if err := srv.ListenAndServe(); err != nil && err != stdhttp.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown error: %v", err)
	}

	log.Println("Server stopped gracefully")
}
