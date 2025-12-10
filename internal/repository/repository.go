package repository

import (
	"context"
	"backend/internal/models"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByID(ctx context.Context, id int64) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
}

type TransactionRepository interface {
	CreateTransaction(ctx context.Context, tx *models.Transaction) error
	GetTransactionByID(ctx context.Context, id int64) (*models.Transaction, error)
	UpdateTransactionStatus(ctx context.Context, id int64, status string) error
}

type BalanceRepository interface {
	GetBalanceByUserID(ctx context.Context, userID int64) (*models.Balance, error)
	UpdateBalance(ctx context.Context, balance *models.Balance) error
	CreateBalance(ctx context.Context, balance *models.Balance) error
}

type AuditRepository interface {
	CreateAuditLog(ctx context.Context, log *models.AuditLog) error
}

type Repository interface {
	UserRepository
	TransactionRepository
	BalanceRepository
	AuditRepository
}
