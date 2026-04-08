package app

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"RESTAPI/internal/config"
	"RESTAPI/internal/handler/httpserver"
	"RESTAPI/internal/repository/postgres"
	"RESTAPI/internal/service"

	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	Server *http.Server
	DB     *pgxpool.Pool
}

func NewApp(ctx context.Context) (*App, error) {
	cfg, err := config.MustLoad()
	if err != nil {
		log.Fatal(err)
	}

	db, err := pgxpool.New(ctx, cfg.DB.DSN)
	if err != nil {
		return nil, fmt.Errorf("db error: %w", err)
	}

	// dependencies
	userRepo := postgres.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := httpserver.NewUserHandler(userService)

	router := httpserver.NewRouter(userHandler)

	server := &http.Server{
		Addr:    cfg.HTTP.Port,
		Handler: router,
	}

	return &App{
		Server: server,
		DB:     db,
	}, nil
}
