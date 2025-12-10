package service

import (
	"context"
	"database/sql"
	"errors"
	"sync"
	"backend/internal/models"
	"backend/internal/repository"
)

type BalanceService struct {
	repo  repository.BalanceRepository
	locks sync.Map // Map[int64]*sync.RWMutex
}

func NewBalanceService(repo repository.BalanceRepository) *BalanceService {
	return &BalanceService{
		repo: repo,
	}
}

func (s *BalanceService) getLock(userID int64) *sync.RWMutex {
	lock, _ := s.locks.LoadOrStore(userID, &sync.RWMutex{})
	return lock.(*sync.RWMutex)
}

func (s *BalanceService) GetBalance(ctx context.Context, userID int64) (*models.Balance, error) {
	mu := s.getLock(userID)
	mu.RLock()
	defer mu.RUnlock()

	bal, err := s.repo.GetBalanceByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &models.Balance{UserID: userID, Amount: 0}, nil
		}
		return nil, err
	}
	return bal, nil
}

func (s *BalanceService) UpdateBalance(ctx context.Context, userID int64, amountDelta int64) error {
	mu := s.getLock(userID)
	mu.Lock()
	defer mu.Unlock()

	balance, err := s.repo.GetBalanceByUserID(ctx, userID)
    
	if err != nil {
		balance = &models.Balance{UserID: userID, Amount: amountDelta}
        return s.repo.CreateBalance(ctx, balance)
	}
	balance.Amount += amountDelta
	return s.repo.UpdateBalance(ctx, balance)
}

func (s *BalanceService) Credit(ctx context.Context, userID int64, amount int64) error {
	if amount <= 0 {
		return errors.New("invalid amount")
	}
	return s.UpdateBalance(ctx, userID, amount)
}

func (s *BalanceService) Debit(ctx context.Context, userID int64, amount int64) error {
	if amount <= 0 {
		return errors.New("invalid amount")
	}

	mu := s.getLock(userID)
	mu.Lock()
	defer mu.Unlock()

	balance, err := s.repo.GetBalanceByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("insufficient funds")
		}
		return err
	}

    if balance.Amount < amount {
        return errors.New("insufficient funds")
    }

	balance.Amount -= amount
	return s.repo.UpdateBalance(ctx, balance)
}
