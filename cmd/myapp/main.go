package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"RESTAPI/internal/app"
)

func main() {
	if err := realMain(); err != nil {
		log.Fatal(err)
	}
}

func realMain() error {
	// initialization context
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// Run
	if err := app.Run(ctx); err != nil {
		return err
	}

	return nil
}
