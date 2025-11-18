package model

import (
	"errors"
	"sync"
	"time"
)

type Balance struct {
	UserID        uint      `gorm:"primaryKey;not null" json:"user_id"`
	Amount        float64   `gorm:"type:decimal(15,2);not null;default:0" json:"amount"`
	LastUpdatedAt time.Time `gorm:"autoUpdateTime" json:"last_updated_at"`

	mu sync.RWMutex `gorm:"-" json:"-"`

	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (b *Balance) Credit(amount float64) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if amount <= 0 {
		return errors.New("credit amount must be greater than zero")
	}

	newAmount := b.Amount + amount
	if newAmount > 10000000 {
		return errors.New("balance would exceed maximum limit of 10,000,000")
	}

	b.Amount = newAmount
	b.LastUpdatedAt = time.Now()
	return nil
}

func (b *Balance) Debit(amount float64) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if amount <= 0 {
		return errors.New("debit amount must be greater than zero")
	}

	if !b.hasSufficientBalanceInternal(amount) {
		return errors.New("insufficient balance")
	}

	b.Amount -= amount
	b.LastUpdatedAt = time.Now()
	return nil
}

func (b *Balance) HasSufficientBalance(amount float64) bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.hasSufficientBalanceInternal(amount)
}

func (b *Balance) hasSufficientBalanceInternal(amount float64) bool {
	return b.Amount >= amount
}
