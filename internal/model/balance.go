package model

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type Balance struct {
	UserID        uint      `gorm:"primaryKey;not null" json:"user_id"`
	Amount        float64   `gorm:"type:decimal(15,2);not null;default:0" json:"amount"`
	LastUpdatedAt time.Time `gorm:"autoUpdateTime" json:"last_updated_at"`

	// Relation
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// BeforeCreate is a GORM hook that runs before creating a record
func (b *Balance) BeforeCreate(tx *gorm.DB) error {
	return b.Validate()
}

// BeforeUpdate is a GORM hook that runs before updating a record
func (b *Balance) BeforeUpdate(tx *gorm.DB) error {
	return b.Validate()
}

// Validate validates all balance fields
func (b *Balance) Validate() error {
	// Validate user ID
	if b.UserID == 0 {
		return errors.New("user_id cannot be empty")
	}

	// Validate amount
	if b.Amount < 0 {
		return errors.New("balance amount cannot be negative")
	}

	if b.Amount > 10000000 {
		return errors.New("balance amount cannot exceed 10,000,000")
	}

	return nil
}

// Credit adds amount to balance
func (b *Balance) Credit(amount float64) error {
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

// Debit subtracts amount from balance
func (b *Balance) Debit(amount float64) error {
	if amount <= 0 {
		return errors.New("debit amount must be greater than zero")
	}

	if !b.HasSufficientBalance(amount) {
		return errors.New("insufficient balance")
	}

	b.Amount -= amount
	b.LastUpdatedAt = time.Now()
	return nil
}

// HasSufficientBalance checks if balance has sufficient amount
func (b *Balance) HasSufficientBalance(amount float64) bool {
	return b.Amount >= amount
}

// IsZero checks if balance is zero
func (b *Balance) IsZero() bool {
	return b.Amount == 0
}

// IsPositive checks if balance is positive
func (b *Balance) IsPositive() bool {
	return b.Amount > 0
}

// TableName specifies the table name for GORM
func (Balance) TableName() string {
	return "balances"
}
