package main

import (
	"go-banking-api/internal/config"
	"go-banking-api/internal/database"
	"go-banking-api/internal/logger"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	cfg := config.Load()
	logger.InitLogger()
	slog.Info("Konfigurasyon ve logger basariyla yuklendi.")

	http.ListenAndServe(cfg.Port, nil)
	slog.Info("Veritabanina baglaniliyor...", "url", cfg.DatabaseURL)
	_, err := database.NewGormConnection(cfg.DatabaseURL)
	if err != nil {
		slog.Error("Veritabani kurulumu basarisiz oldu", "hata", err)
		os.Exit(1)
	}

	slog.Info("Veritabani baglantisi ve tablo migrasyonu basariyla tamamlandi.")
	slog.Info("Tablolarin olustugunu veritabani yonetim aracinizdan (DBeaver, pgAdmin) kontrol edebilirsiniz.")
}
