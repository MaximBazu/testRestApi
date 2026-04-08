package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"RESTAPI/internal/app"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	application, err := app.NewApp(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer application.DB.Close()

	go func() {
		log.Println("server started on", application.Server.Addr)

		if err := application.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	// ждём сигнал
	<-ctx.Done()
	log.Println("shutting down server...")

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := application.Server.Shutdown(ctxShutdown); err != nil {
		log.Fatal("server shutdown failed:", err)
	}

	log.Println("server exited properly")
}
