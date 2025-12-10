package models

import (
	"errors"
	"strings"
	"time"
)

const (
	RoleUser  = "user"
	RoleAdmin = "admin"
)
const (
	TxTypeDeposit  = "deposit"
	TxTypeWithdraw = "withdraw"
	TxTypeTransfer = "transfer"
)
const (
	TxStatusPending   = "pending"
	TxStatusCompleted = "completed"
	TxStatusFailed    = "failed"
)
type User struct {
	ID           int64     `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (u *User) Validate() error {
	if len(u.Username) < 3 {
		return errors.New("username must be at least 3 characters")
	}
	if !strings.Contains(u.Email, "@") {
		return errors.New("invalid email format")
	}
	if len(u.PasswordHash) == 0 {
		return errors.New("password is required")
	}
	return nil
}

type Transaction struct {
	ID         int64     `json:"id"`
	FromUserID *int64    `json:"from_user_id,omitempty"` // Nullable for deposits
	ToUserID   *int64    `json:"to_user_id,omitempty"`   // Nullable for withdrawals (if applicable)
	Amount     int64     `json:"amount"`                 // In cents
	Type       string    `json:"type"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}

type Balance struct {
	UserID        int64     `json:"user_id"`
	Amount        int64     `json:"amount"` // In cents
	LastUpdatedAt time.Time `json:"last_updated_at"`
}

type AuditLog struct {
	ID         int64     `json:"id"`
	EntityType string    `json:"entity_type"`
	EntityID   int64     `json:"entity_id"`
	Action     string    `json:"action"`
	Details    string    `json:"details"`
	CreatedAt  time.Time `json:"created_at"`
}
