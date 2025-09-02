package main

import (
	"auth-service/src/config"
	"auth-service/src/repository"
	"auth-service/src/server"
	"auth-service/src/service"
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found, using system environment variables.")
	}
	cfg := config.Load()

	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}
	defer pool.Close()
	log.Println("Successfully connected to PostgreSQL.")

	userRepo := repository.NewUser(pool)
	userService := service.NewUserService(userRepo, cfg.JWTSecret)
	httpServer := server.NewServer(cfg, userService)

	httpServer.Run()
}
