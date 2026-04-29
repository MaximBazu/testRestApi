package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"RESTAPI/internal/config"
	"RESTAPI/internal/handler/httpserver"
	"RESTAPI/internal/repository/postgres"
	"RESTAPI/internal/service"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Run(ctx context.Context) error {
	// канал для ошибок сервера
	serverErr := make(chan error, 1)

	// config
	cfg, err := config.MustLoad()
	if err != nil {
		return err
	}

	// database (pool)
	db, err := pgxpool.New(ctx, cfg.DB.DSN)
	if err != nil {
		return fmt.Errorf("db error: %w", err)
	}
	defer db.Close()

	pingCtx, pingCancel := context.WithTimeout(ctx, 3*time.Second)
	defer pingCancel()

	if err := db.Ping(pingCtx); err != nil {
		return fmt.Errorf("db ping error: %w", err)
	}

	// dependencies
	userRepo := postgres.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := httpserver.NewUserHandler(userService)

	productRepo := postgres.NewProductRepository(db)
	productService := service.NewProductService(productRepo)
	productHandler := httpserver.NewProductHandler(productService)

	orderRepo := postgres.NewOrderRepository(db)
	orderService := service.NewOrderService(orderRepo)
	orderHandler := httpserver.NewOrderHandler(orderService)

	orderItemRepo := postgres.NewOrderItemRepository(db)
	orderItemService := service.NewOrderItemService(orderItemRepo)
	orderItemHandler := httpserver.NewOrderItemHandler(orderItemService)

	productSizeRepo := postgres.NewProductSizeRepository(db)
	productSizeService := service.NewProductSizeService(productSizeRepo)
	productSizeHandler := httpserver.NewProductSizeHandler(productSizeService)

	productImageRepo := postgres.NewProductImageRepository(db)
	productImageService := service.NewProductImageService(productImageRepo)
	productImageHandler := httpserver.NewProductImageHandler(productImageService)

	router := httpserver.NewRouter(
		userHandler,
		productHandler,
		orderHandler,
		orderItemHandler,
		productSizeHandler,
		productImageHandler,
	)
	server := &http.Server{
		Addr:    cfg.HTTP.Port,
		Handler: router,
	}

	// start server
	go func() {
		log.Println("server started on", server.Addr)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr <- fmt.Errorf("listen error: %w", err)
		}
		close(serverErr)
	}()

	// ждём либо сигнал остановки, либо ошибку сервера
	select {
	case <-ctx.Done(): // graceful shutdown
		log.Println("shutting down server...")
	case err := <-serverErr:
		return fmt.Errorf("server error: %w", err)
	}

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctxShutdown); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	log.Println("server exited properly")
	return nil
}
