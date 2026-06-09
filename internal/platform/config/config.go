package config

import (
	"os"
)

type Config struct {
	Env       string
	Port      string
	DBURL     string
	JWTSecret string
}

func Load() *Config {
	return &Config{
		Env:       getEnv("APP_ENV", "development"),
		Port:      getEnv("PORT", "8080"),
		DBURL:     getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/core?sslmode=disable"),
		JWTSecret: getEnv("JWT_SECRET", "chave-secreta-padrao-apenas-para-desenvolvimento"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}