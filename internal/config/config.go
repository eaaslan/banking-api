package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port         string
	DBHost       string
	DBPort       string
	DBUser       string
	DBPassword   string
	DBName       string
	Environment  string
	LogLevel     string
	AuthSecret    string
	RedisHost     string
	RedisPort     string
	RedisPassword string
	OTLPEndpoint  string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	return &Config{
		Port:          getEnv("PORT", "8080"),
		DBHost:        getEnv("DB_HOST", "localhost"),
		DBPort:        getEnv("DB_PORT", "5432"),
		DBUser:        getEnv("DB_USER", "postgres"),
		DBPassword:    getEnv("DB_PASSWORD", "postgres"),
		DBName:        getEnv("DB_NAME", "bank"),
		Environment:   getEnv("ENV", "development"),
		LogLevel:      getEnv("LOG_LEVEL", "info"),
		AuthSecret:    getEnv("AUTH_SECRET", "super-secret-key"),
		RedisHost:     getEnv("REDIS_HOST", "localhost"),
		RedisPort:     getEnv("REDIS_PORT", "6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		OTLPEndpoint:  getEnv("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4318"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
