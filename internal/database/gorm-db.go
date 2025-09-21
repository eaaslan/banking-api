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
		return nil, fmt.Errorf("veritabanina baglanilamadi: %w", err)
	}

	slog.Info("Veritabani migrasyonu baslatiliyor...")
	err = db.AutoMigrate(
		&model.User{},
		&model.Transaction{},
		&model.Balance{},
		&model.AuditLog{},
	)
	if err != nil {
		return nil, fmt.Errorf("veritabani migrasyonu basarisiz oldu: %w", err)
	}

	slog.Info("Veritabani migrasyonu basariyla tamamlandi.")
	return db, nil
}
