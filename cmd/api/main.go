package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"backend/internal/config"
	"backend/internal/db"
	apiHandler "backend/internal/handler"
	"backend/internal/middleware"
	"backend/internal/repository"
	"backend/internal/router"
	"backend/internal/service"
	"backend/internal/worker"
)

func main() {
	cfg := config.Load()
	var handler slog.Handler
	if cfg.Environment == "production" {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	} else {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
	}
	logger := slog.New(handler)
	slog.SetDefault(logger)

	logger.Info("Initializing application...", "env", cfg.Environment)

	database, err := db.NewDB(cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	if err != nil {
		logger.Error("Failed to initialize database", "error", err)
		os.Exit(1)
	}
	defer database.Close()

	if err := db.RunMigrations(database, "migrations"); err != nil {
		logger.Error("Failed to run migrations", "error", err)
		os.Exit(1)
	}
	logger.Info("Database initialized and migrations run")

	repo := repository.NewPostgresRepository(database)
	userSvc := service.NewUserService(repo, cfg.AuthSecret)
	balSvc := service.NewBalanceService(repo)
	txSvc := service.NewTransactionService(repo, balSvc)
	poolCtx, poolCancel := context.WithCancel(context.Background())
	defer poolCancel()

	pool := worker.NewPool(5, 100, txSvc.ProcessTransaction)
	pool.Start(poolCtx)
	txSvc.SetPool(pool)

	h := apiHandler.NewHandler(userSvc, txSvc, balSvc)

	r := router.NewRouter()
	r.Use(middleware.Logger, middleware.Recovery, middleware.CORS, middleware.RateLimit)

	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	
	// Public routes
	r.HandleFunc("/api/v1/auth/register", h.Register)
	r.HandleFunc("/api/v1/auth/login", h.Login)
	r.HandleFunc("/api/v1/auth/refresh", h.Refresh)

	// Protected routes
	authMw := middleware.Auth(userSvc)
	
	// Transaction Routes
	r.HandleFunc("/api/v1/transactions", h.CreateTransaction, authMw)
	r.HandleFunc("/api/v1/transactions/history", h.GetTransactionHistory, authMw)
	
	// Balance Routes
	r.HandleFunc("/api/v1/balances/current", h.GetBalance, authMw)
	r.HandleFunc("/api/v1/balances/historical", h.GetBalanceHistory, authMw)
	
	// User Routes
	roleMw := middleware.Role("admin")
	r.HandleFunc("/api/v1/users", h.ListUsers, authMw, roleMw)
	r.HandleFunc("/api/v1/users/delete", h.DeleteUser, authMw, roleMw) // Using query param ?id=

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	go func() {
		logger.Info("Server listening", "port", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Server startup failed", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	sign := <-quit
	logger.Info("Shutdown signal received", "signal", sign.String())
	
	poolCancel()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", "error", err)
	}

	logger.Info("Server exited gracefully")
}
