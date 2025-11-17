package database

import (
	"fmt"
	"go-banking-api/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log/slog"
)

func NewGormConnection(databaseURL string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("Cant connect to database: %w", err)
	}

	slog.Info("Database migrations starting ...")
	err = db.AutoMigrate(&model.User{}, &model.Balance{}, &model.Transaction{}, &model.AuditLog{})
	if err != nil {
		return nil, fmt.Errorf("Database migrations failed: %w", err)
	}

	slog.Info("Database migration is successful.")
	return db, nil
}
