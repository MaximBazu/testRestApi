package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTP HTTPConfig
	DB   DBConfig
}

type HTTPConfig struct {
	Port string
}

type DBConfig struct {
	DSN string
}

func MustLoad() (*Config, error) {
	godotenv.Load()

	cfg := &Config{
		HTTP: HTTPConfig{
			Port: getEnv("HTTP_PORT", ":8080"),
		},
		DB: DBConfig{
			DSN: getEnv("DB_DSN", ""),
		},
	}

	if cfg.DB.DSN == "" {
		return nil, fmt.Errorf("DB_DSN is required")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
