package logger

import (
	"log/slog"
	"os"
)

func InitLogger() {
	var logLevel slog.Level

	logLevel = slog.LevelInfo

	//if os.Getenv("GO_ENV") == "development" {
	//	logLevel = slog.LevelDebug
	//}

	jsonHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	})

	slog.SetDefault(slog.New(jsonHandler))
}
