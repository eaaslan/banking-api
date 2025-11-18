package domain

import (
	"context"
	"go-banking-api/internal/model"
)

type TransactionService interface {
	CreateTransaction(ctx context.Context, transaction *model.Transaction) error
	GetTransactionByID(ctx context.Context, id uint) (*model.Transaction, error)
	GetUserTransactions(ctx context.Context, userID uint, limit, offset int) ([]*model.Transaction, error)
	ProcessTransaction(ctx context.Context, transactionID uint) error
	CancelTransaction(ctx context.Context, transactionID uint) error
}

// In internal/domain/balance_service.go

type BalanceService interface {
	GetBalance(ctx context.Context, userID uint) (*model.Balance, error)
	CreditBalance(ctx context.Context, userID uint, amount float64) error
	DebitBalance(ctx context.Context, userID uint, amount float64) error
	TransferBalance(ctx context.Context, fromUserID, toUserID uint, amount float64) error
}
