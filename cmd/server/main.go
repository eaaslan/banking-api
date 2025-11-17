package main

import (
	"context"
	"errors"
	"fmt"
	"go-banking-api/internal/config"
	"go-banking-api/internal/database"
	"go-banking-api/internal/handler"
	"go-banking-api/internal/repository"
	"go-banking-api/internal/router"
	"go-banking-api/internal/service"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg := config.Load()
	setupLogger(cfg.Environment)
	slog.Info("Starting application...")

	db, err := database.NewGormConnection(cfg.DatabaseURL)
	if err != nil {
		slog.Error("Database connection failed", "error", err)
		os.Exit(1)
	}
	slog.Info("Database connected successfully")

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	routes := router.SetupRoutes(userHandler)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      routes,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		slog.Info("Server starting", "address", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Server failed to start", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	recvdSignal := <-quit
	slog.Info("Shutdown signal received", "signal", recvdSignal.String())

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("Server shutdown failed", "error", err)
		os.Exit(1)
	}
	slog.Info("Server stopped gracefully")
}

func setupLogger(environment string) {
	var logger *slog.Logger
	if environment == "development" {
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	} else {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	}
	slog.SetDefault(logger)
}
