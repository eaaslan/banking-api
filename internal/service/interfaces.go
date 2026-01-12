package service

import (
	"context"
	"backend/internal/models"
)

type UserServiceInterface interface {
	Register(ctx context.Context, username, email, password string) (*models.User, error)
	Authenticate(ctx context.Context, email, password string) (*models.User, error)
	// Add other methods as needed
}

type TransactionServiceInterface interface {
	Create(ctx context.Context, fromID, toID *int64, amount int64, typeStr string) (*models.Transaction, error)
	ProcessTransaction(ctx context.Context, tx *models.Transaction) error
	// GetHistory(ctx context.Context, userID int64) ([]*models.Transaction, error) // To be implemented
}

type BalanceServiceInterface interface {
	GetBalance(ctx context.Context, userID int64) (*models.Balance, error)
	UpdateBalance(ctx context.Context, userID int64, amountDelta int64) error
	Credit(ctx context.Context, userID int64, amount int64) error
	Debit(ctx context.Context, userID int64, amount int64) error
}
