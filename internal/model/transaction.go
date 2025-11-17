package model

import (
	"errors"
	"strings"

	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	FromUserID *uint   `gorm:"index" json:"from_user_id"` // Sender (nullable - for credit operations)
	ToUserID   *uint   `gorm:"index" json:"to_user_id"`   // Receiver (nullable - for debit operations)
	Amount     float64 `gorm:"type:decimal(15,2);not null" json:"amount"`
	Type       string  `gorm:"type:varchar(20);not null" json:"type"`                     // "credit", "debit", "transfer"
	Status     string  `gorm:"type:varchar(20);not null;default:'pending'" json:"status"` // "pending", "completed", "failed"

	// Relations
	FromUser *User `gorm:"foreignKey:FromUserID" json:"from_user,omitempty"`
	ToUser   *User `gorm:"foreignKey:ToUserID" json:"to_user,omitempty"`
}

// Transaction types
const (
	TransactionTypeCredit   = "credit"   // Deposit money
	TransactionTypeDebit    = "debit"    // Withdraw money
	TransactionTypeTransfer = "transfer" // Transfer between accounts
)

// Transaction statuses
const (
	TransactionStatusPending   = "pending"
	TransactionStatusCompleted = "completed"
	TransactionStatusFailed    = "failed"
)

// BeforeCreate is a GORM hook that runs before creating a record
func (t *Transaction) BeforeCreate(tx *gorm.DB) error {
	return t.Validate()
}

// BeforeUpdate is a GORM hook that runs before updating a record
func (t *Transaction) BeforeUpdate(tx *gorm.DB) error {
	return t.Validate()
}

// Validate validates all transaction fields
func (t *Transaction) Validate() error {
	// Validate amount
	if err := t.validateAmount(); err != nil {
		return err
	}

	// Validate type
	if err := t.validateType(); err != nil {
		return err
	}

	// Validate status
	if err := t.validateStatus(); err != nil {
		return err
	}

	// Validate user IDs based on transaction type
	if err := t.validateUserIDs(); err != nil {
		return err
	}

	return nil
}

// validateAmount validates transaction amount
func (t *Transaction) validateAmount() error {
	if t.Amount <= 0 {
		return errors.New("amount must be greater than zero")
	}

	if t.Amount > 1000000 {
		return errors.New("amount cannot exceed 1,000,000")
	}

	return nil
}

// validateType validates transaction type
func (t *Transaction) validateType() error {
	t.Type = strings.TrimSpace(strings.ToLower(t.Type))

	if t.Type == "" {
		return errors.New("transaction type cannot be empty")
	}

	validTypes := []string{TransactionTypeCredit, TransactionTypeDebit, TransactionTypeTransfer}
	for _, vt := range validTypes {
		if t.Type == vt {
			return nil
		}
	}

	return errors.New("transaction type must be 'credit', 'debit', or 'transfer'")
}

// validateStatus validates transaction status
func (t *Transaction) validateStatus() error {
	t.Status = strings.TrimSpace(strings.ToLower(t.Status))

	// If status is empty, set default
	if t.Status == "" {
		t.Status = TransactionStatusPending
		return nil
	}

	validStatuses := []string{TransactionStatusPending, TransactionStatusCompleted, TransactionStatusFailed}
	for _, vs := range validStatuses {
		if t.Status == vs {
			return nil
		}
	}

	return errors.New("status must be 'pending', 'completed', or 'failed'")
}

// validateUserIDs validates user IDs based on transaction type
func (t *Transaction) validateUserIDs() error {
	switch t.Type {
	case TransactionTypeCredit:
		// Credit: ToUserID required, FromUserID must be nil
		if t.ToUserID == nil || *t.ToUserID == 0 {
			return errors.New("to_user_id is required for credit transactions")
		}
		if t.FromUserID != nil && *t.FromUserID != 0 {
			return errors.New("from_user_id must be empty for credit transactions")
		}

	case TransactionTypeDebit:
		// Debit: FromUserID required, ToUserID must be nil
		if t.FromUserID == nil || *t.FromUserID == 0 {
			return errors.New("from_user_id is required for debit transactions")
		}
		if t.ToUserID != nil && *t.ToUserID != 0 {
			return errors.New("to_user_id must be empty for debit transactions")
		}

	case TransactionTypeTransfer:
		// Transfer: Both required and must be different
		if t.FromUserID == nil || *t.FromUserID == 0 {
			return errors.New("from_user_id is required for transfer transactions")
		}
		if t.ToUserID == nil || *t.ToUserID == 0 {
			return errors.New("to_user_id is required for transfer transactions")
		}
		if *t.FromUserID == *t.ToUserID {
			return errors.New("cannot transfer to the same account")
		}
	}

	return nil
}

// IsPending checks if transaction is pending
func (t *Transaction) IsPending() bool {
	return t.Status == TransactionStatusPending
}

// IsCompleted checks if transaction is completed
func (t *Transaction) IsCompleted() bool {
	return t.Status == TransactionStatusCompleted
}

// IsFailed checks if transaction is failed
func (t *Transaction) IsFailed() bool {
	return t.Status == TransactionStatusFailed
}

// MarkCompleted marks transaction as completed
func (t *Transaction) MarkCompleted() {
	t.Status = TransactionStatusCompleted
}

// MarkFailed marks transaction as failed
func (t *Transaction) MarkFailed() {
	t.Status = TransactionStatusFailed
}

// IsCredit checks if transaction is a credit operation
func (t *Transaction) IsCredit() bool {
	return t.Type == TransactionTypeCredit
}

// IsDebit checks if transaction is a debit operation
func (t *Transaction) IsDebit() bool {
	return t.Type == TransactionTypeDebit
}

// IsTransfer checks if transaction is a transfer operation
func (t *Transaction) IsTransfer() bool {
	return t.Type == TransactionTypeTransfer
}

// TableName specifies the table name for GORM
func (Transaction) TableName() string {
	return "transactions"
}
