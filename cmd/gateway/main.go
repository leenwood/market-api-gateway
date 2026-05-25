// Package main is the entry point for the API Gateway service.
//
//	@title			Market API Gateway
//	@version		1.0
//	@description	Single entry point for all marketplace clients (web, mobile). Routes traffic between upstream services, validates JWT, enforces rate limits and CORS.
//
//	@contact.name	Platform Team
//
//	@license.name	MIT
//
//	@host		localhost:8080
//	@BasePath	/
//
//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				JWT access token. Format: "Bearer <token>"
package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	_ "market-api-gateway/docs" // swagger generated docs
	"market-api-gateway/internal/app/service"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := service.RunServer(ctx); err != nil {
		log.Fatalf("gateway: %v", err)
	}
}
