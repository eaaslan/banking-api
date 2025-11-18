package domain

import (
	"context"
	"go-banking-api/internal/model"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error

	GetByID(ctx context.Context, id uint) (*model.User, error)

	GetByUsername(ctx context.Context, username string) (*model.User, error)

	GetByEmail(ctx context.Context, email string) (*model.User, error)

	GetAll(ctx context.Context, limit, offset int) ([]*model.User, error)

	Update(ctx context.Context, user *model.User) error

	Delete(ctx context.Context, id uint) error
}

type TransactionRepository interface {
	Create(ctx context.Context, transaction *model.Transaction) error
	GetByID(ctx context.Context, id uint) (*model.Transaction, error)
	GetByUserID(ctx context.Context, userID uint, limit, offset int) ([]*model.Transaction, error)
	GetPendingTransactions(ctx context.Context, limit, offset int) ([]*model.Transaction, error)
	Update(ctx context.Context, transaction *model.Transaction) error
	Delete(ctx context.Context, id uint) error
}

type BalanceRepository interface {
	GetByUserID(ctx context.Context, userID uint) (*model.Balance, error)
	Create(ctx context.Context, balance *model.Balance) error
	Update(ctx context.Context, balance *model.Balance) error
}
