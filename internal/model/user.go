package model

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username     string `gorm:"unique;not null;size:50"`
	Email        string `gorm:"unique;not null;size:255"`
	PasswordHash string `gorm:"not null"`
}

type Transaction struct {
	gorm.Model
	FromUserID uint    `gorm:"index"`
	ToUserID   uint    `gorm:"index"`
	Amount     float64 `gorm:"not null;type:decimal(10,2)"`
	Type       string  `gorm:"not null;size:20"`
	Status     string  `gorm:"default:PENDING;size:20"`
}

type Balance struct {
	gorm.Model
	UserID uint    `gorm:"uniqueIndex;not null"`
	Amount float64 `gorm:"default:0;type:decimal(10,2)"`
}

type AuditLog struct {
	gorm.Model
	EntityID uint   `gorm:"not null"`
	Action   string `gorm:"not null;size:20"`
	Details  string `gorm:"type:text"`
}
