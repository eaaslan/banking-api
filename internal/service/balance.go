package service

import (
	"context"
	"database/sql"
	"errors"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"backend/internal/cache"
	"backend/internal/models"
	"backend/internal/repository"
)

type BalanceService struct {
	repo  repository.Repository
	locks sync.Map
	redis *cache.RedisClient
}

func NewBalanceService(repo repository.Repository, redisClient *cache.RedisClient) *BalanceService {
	return &BalanceService{
		repo:  repo,
		redis: redisClient,
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

	// Try cache
	key := fmt.Sprintf("balance:%d", userID)
	val, err := s.redis.Client.Get(ctx, key).Result()
	if err == nil {
		var bal models.Balance
		if err := json.Unmarshal([]byte(val), &bal); err == nil {
			return &bal, nil
		}
	}

	bal, err := s.repo.GetBalanceByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &models.Balance{UserID: userID, Amount: 0}, nil
		}
		return nil, err
	}

	// Set cache
	if bytes, err := json.Marshal(bal); err == nil {
		s.redis.Client.Set(ctx, key, bytes, 10*time.Minute)
	}

	return bal, nil
}

func (s *BalanceService) GetHistory(ctx context.Context, userID int64) ([]*models.AuditLog, error) {
    return s.repo.GetAuditLogsByEntity(ctx, "user", userID)
}

func (s *BalanceService) UpdateBalance(ctx context.Context, userID int64, amountDelta int64) error {
	mu := s.getLock(userID)
	mu.Lock()
	defer mu.Unlock()

	balance, err := s.repo.GetBalanceByUserID(ctx, userID)
	if err != nil {
		balance = &models.Balance{UserID: userID, Amount: amountDelta}
		if err := s.repo.CreateBalance(ctx, balance); err != nil {
            return err
        }
	} else {
	    balance.Amount += amountDelta
	    if err := s.repo.UpdateBalance(ctx, balance); err != nil {
            return err
        }
    }
    
	// Invalidate cache
	s.redis.Client.Del(ctx, fmt.Sprintf("balance:%d", userID))

    // Audit Log
    _ = s.repo.CreateAuditLog(ctx, &models.AuditLog{
        EntityType: "user",
        EntityID:   userID,
        Action:     "balance_update",
        Details:    fmt.Sprintf("amount_delta: %d", amountDelta),
    })
    
	return nil
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
	if err := s.repo.UpdateBalance(ctx, balance); err != nil {
		return err
	}
	
	// Invalidate cache
	s.redis.Client.Del(ctx, fmt.Sprintf("balance:%d", userID))
	
	return nil
}
