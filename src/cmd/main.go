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
		log.Println("Arquivo .env não encontrado, usando variáveis de ambiente do sistema.")
	}
	cfg := config.Load()

	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Falha ao conectar com o PostgreSQL: %v", err)
	}
	defer pool.Close()
	log.Println("Conectado ao PostgreSQL com sucesso.")

	userRepo := repository.NewPostgresUserRepository(pool)
	userService := service.NewUserService(userRepo, cfg.JWTSecret)
	httpServer := server.NewServer(cfg, userService)

	httpServer.Run()
}
