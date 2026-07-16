package config

import (
	"fmt"
	"os"
)

type Config struct {
	AppPort     string
	DatabaseURL string
	RedisAddr   string
}

func Load() (*Config, error) {
	required := []string{"APP_PORT", "DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME", "DB_SSLMODE", "REDIS_ADDR"}
	for _, key := range required {
		if os.Getenv(key) == "" {
			return nil, fmt.Errorf("missing required environment variable %s", key)
		}
	}

	databaseURL := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"),
	)

	return &Config{
		AppPort:     os.Getenv("APP_PORT"),
		DatabaseURL: databaseURL,
		RedisAddr:   os.Getenv("REDIS_ADDR"),
	}, nil
}
