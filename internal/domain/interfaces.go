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
