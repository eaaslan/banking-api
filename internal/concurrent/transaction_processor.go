package concurrent

import (
	"context"
	"errors"
	"go-banking-api/internal/domain"
	"go-banking-api/internal/model"
	"log"
)

type DefaultTransactionProcessor struct {
	transactionRepo domain.TransactionRepository
	balanceRepo     domain.BalanceRepository
}

func NewTransactionProcessor(
	transactionRepo domain.TransactionRepository,
	balanceRepo domain.BalanceRepository,
) *DefaultTransactionProcessor {
	return &DefaultTransactionProcessor{
		transactionRepo: transactionRepo,
		balanceRepo:     balanceRepo,
	}
}

func (p *DefaultTransactionProcessor) Process(ctx context.Context, transaction *model.Transaction) error {
	if !transaction.IsPending() {
		return nil
	}

	dbTransaction, err := p.transactionRepo.GetByID(ctx, transaction.ID)
	if err != nil {
		return err
	}

	if !dbTransaction.IsPending() {
		return nil
	}

	switch {
	case dbTransaction.IsCredit():
		err = p.processCredit(ctx, dbTransaction)
	case dbTransaction.IsDebit():
		err = p.processDebit(ctx, dbTransaction)
	case dbTransaction.IsTransfer():
		err = p.processTransfer(ctx, dbTransaction)
	default:
		err = errors.New("unknown transaction type")
	}

	if err != nil {
		dbTransaction.MarkFailed()
		log.Printf("Transaction %d failed: %v", dbTransaction.ID, err)
	} else {
		dbTransaction.MarkCompleted()
		log.Printf("Transaction %d completed successfully", dbTransaction.ID)
	}

	updateErr := p.transactionRepo.Update(ctx, dbTransaction)
	if updateErr != nil {
		log.Printf("Failed to update transaction %d status: %v", dbTransaction.ID, updateErr)
		return updateErr
	}

	return err
}

func (p *DefaultTransactionProcessor) processCredit(ctx context.Context, transaction *model.Transaction) error {
	if transaction.ToUserID == nil {
		return errors.New("to_user_id is required for credit transactions")
	}

	balance, err := p.balanceRepo.GetByUserID(ctx, *transaction.ToUserID)
	if err != nil {
		return err
	}

	if err := balance.Credit(transaction.Amount); err != nil {
		return err
	}

	return p.balanceRepo.Update(ctx, balance)
}

func (p *DefaultTransactionProcessor) processDebit(ctx context.Context, transaction *model.Transaction) error {
	if transaction.FromUserID == nil {
		return errors.New("from_user_id is required for debit transactions")
	}

	balance, err := p.balanceRepo.GetByUserID(ctx, *transaction.FromUserID)
	if err != nil {
		return err
	}

	if !balance.HasSufficientBalance(transaction.Amount) {
		return errors.New("insufficient balance")
	}

	if err := balance.Debit(transaction.Amount); err != nil {
		return err
	}

	return p.balanceRepo.Update(ctx, balance)
}

func (p *DefaultTransactionProcessor) processTransfer(ctx context.Context, transaction *model.Transaction) error {
	if transaction.FromUserID == nil {
		return errors.New("from_user_id is required for transfer transactions")
	}
	if transaction.ToUserID == nil {
		return errors.New("to_user_id is required for transfer transactions")
	}

	fromBalance, err := p.balanceRepo.GetByUserID(ctx, *transaction.FromUserID)
	if err != nil {
		return err
	}

	if !fromBalance.HasSufficientBalance(transaction.Amount) {
		return errors.New("insufficient balance")
	}

	toBalance, err := p.balanceRepo.GetByUserID(ctx, *transaction.ToUserID)
	if err != nil {
		return err
	}

	if err := fromBalance.Debit(transaction.Amount); err != nil {
		return err
	}

	if err := toBalance.Credit(transaction.Amount); err != nil {
		_ = fromBalance.Credit(transaction.Amount)
		return err
	}

	if err := p.balanceRepo.Update(ctx, fromBalance); err != nil {
		_ = fromBalance.Credit(transaction.Amount)
		_ = toBalance.Debit(transaction.Amount)
		return err
	}

	if err := p.balanceRepo.Update(ctx, toBalance); err != nil {
		_ = fromBalance.Credit(transaction.Amount)
		_ = p.balanceRepo.Update(ctx, fromBalance)
		return err
	}

	return nil
}
