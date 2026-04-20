// internal/config/config.go
package config

import (
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	HTTP struct {
		Port        string        `envconfig:"HTTP_PORT" default:":8080"`
		ReadTimeout time.Duration `envconfig:"HTTP_READ_TIMEOUT" default:"15s"`
	}
	DB struct {
		DSN          string `envconfig:"DB_DSN" required:"true"`
		MaxOpenConns int    `envconfig:"DB_MAX_OPEN_CONNS" default:"25"`
		MaxIdleConns int    `envconfig:"DB_MAX_IDLE_CONNS" default:"5"`
	}
	AssetsDir string `envconfig:"ASSETS_DIR" default:"./assets"`
}

func MustLoad() (*Config, error) {
	_ = godotenv.Load()
	var cfg Config
	// Префикс "APP" → переменные вида APP_DB_DSN, APP_HTTP_PORT
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
