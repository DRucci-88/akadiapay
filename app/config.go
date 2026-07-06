package app

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	APP_PORT   int
	APP_ENV    string
	DB_DSN     string
	JWT_SECRET string
}

func LoadConfig() *AppConfig {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	var envFile string
	if env == "production" {
		envFile = ".env.production"
	} else {
		envFile = ".env.development"
	}

	if err := godotenv.Load(envFile); err != nil {
		log.Println("Peringatan: File .env tidak ditemukan, sistem akan menggunakan env bawaan OS")
	}

	appPort, _ := strconv.Atoi(getEnv("APP_PORT", "8080"))
	defaultDSN := "postgres://postgres:12345678@localhost:5432/akadia?sslmode=disable"

	return &AppConfig{
		APP_PORT:   appPort,
		APP_ENV:    getEnv("APP_ENV", "development"),
		DB_DSN:     getEnv("DB_DSN", defaultDSN),
		JWT_SECRET: getEnv("JWT_SECRET", "KunciRahasiaNegaraSangatRahasiaSekali"),
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
