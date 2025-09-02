package seeder

import (
	"auth-service/src/domain"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TestSeeder struct {
	db *pgxpool.Pool
}

func NewTestSeeder(db *pgxpool.Pool) *TestSeeder {
	return &TestSeeder{db: db}
}

func (s *TestSeeder) InsertUser(ctx context.Context, user *domain.User) error {
	query := `INSERT INTO users (id, name, email, password_hash, created_at) VALUES ($1, $2, $3, $4, $5)`
	_, err := s.db.Exec(ctx, query, user.ID, user.Name, user.Email, user.PasswordHash, user.CreatedAt)
	return err
}

func (s *TestSeeder) TruncateTables(ctx context.Context) error {
	_, err := s.db.Exec(ctx, "TRUNCATE TABLE users RESTART IDENTITY")
	return err
}
