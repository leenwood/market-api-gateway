package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"market-api-gateway/internal/app/service"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := service.RunServer(ctx); err != nil {
		log.Fatalf("gateway: %v", err)
	}
}
