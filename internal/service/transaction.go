package service

import (
	"context"
	"errors"
	"backend/internal/models"
	"backend/internal/repository"
	"backend/internal/worker"
)

type TransactionService struct {
	repo       repository.TransactionRepository
	balanceSvc *BalanceService
	pool       *worker.Pool
}

func NewTransactionService(repo repository.TransactionRepository, balanceSvc *BalanceService) *TransactionService {
	return &TransactionService{
		repo:       repo,
		balanceSvc: balanceSvc,
	}
}

func (s *TransactionService) SetPool(pool *worker.Pool) {
	s.pool = pool
}

func (s *TransactionService) Create(ctx context.Context, fromID, toID *int64, amount int64, typeStr string) (*models.Transaction, error) {
	tx := &models.Transaction{
		FromUserID: fromID,
		ToUserID:   toID,
		Amount:     amount,
		Type:       typeStr,
		Status:     models.TxStatusPending,
	}

	if err := s.repo.CreateTransaction(ctx, tx); err != nil {
		return nil, err
	}

	if s.pool != nil {
		s.pool.Submit(tx)
	} else {
		return nil, errors.New("worker pool not initialized")
	}

	return tx, nil
}

func (s *TransactionService) ProcessTransaction(ctx context.Context, tx *models.Transaction) error {
	var err error

	switch tx.Type {
	case models.TxTypeDeposit:
		if tx.ToUserID == nil {
			err = errors.New("missing to_user")
		} else {
			err = s.balanceSvc.Credit(ctx, *tx.ToUserID, tx.Amount)
		}

	case models.TxTypeWithdraw:
		if tx.FromUserID == nil {
			err = errors.New("missing from_user")
		} else {
			err = s.balanceSvc.Debit(ctx, *tx.FromUserID, tx.Amount)
		}

	case models.TxTypeTransfer:
		if tx.FromUserID == nil || tx.ToUserID == nil {
			err = errors.New("invalid transfer users")
		} else {
			err = s.balanceSvc.Debit(ctx, *tx.FromUserID, tx.Amount)
			if err == nil {
				if creditErr := s.balanceSvc.Credit(ctx, *tx.ToUserID, tx.Amount); creditErr != nil {
					_ = s.balanceSvc.Credit(ctx, *tx.FromUserID, tx.Amount)
					err = creditErr
				}
			}
		}

	default:
		err = errors.New("unknown transaction type")
	}

	status := models.TxStatusCompleted
	if err != nil {
		status = models.TxStatusFailed
	}
	_ = s.repo.UpdateTransactionStatus(ctx, tx.ID, status)

	return err
}
