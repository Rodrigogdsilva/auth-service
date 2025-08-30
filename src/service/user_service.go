package service

import (
	"auth-service/src/domain"
	"auth-service/src/jwt"
	"auth-service/src/repository"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Register(ctx context.Context, name, email, password string) (*domain.User, error)
	Login(ctx context.Context, email, password string) (string, error)
	GetProfile(ctx context.Context, userID string) (*domain.User, error)
	ValidateToken(tokenString string) (map[string]interface{}, error)
}

type userService struct {
	repo      repository.UserRepository
	jwtSecret string
}

func NewUserService(repo repository.UserRepository, jwtSecret string) UserService {
	return &userService{repo: repo, jwtSecret: jwtSecret}
}

func (s *userService) Register(ctx context.Context, name, email, password string) (*domain.User, error) {
	if name == "" || email == "" || password == "" {
		return nil, errors.New("nome, email e senha são obrigatórios")
	}
	if len(password) < 8 {
		return nil, errors.New("senha deve ter no mínimo 8 caracteres")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("falha ao criptografar senha")
	}

	user := &domain.User{
		ID:           uuid.NewString(),
		Name:         name,
		Email:        email,
		PasswordHash: string(hashedPassword),
		CreatedAt:    time.Now().UTC(),
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return "", errors.New("credenciais inválidas")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("credenciais inválidas")
	}
	return jwt.CreateToken(user, s.jwtSecret)
}

func (s *userService) GetProfile(ctx context.Context, userID string) (*domain.User, error) {
	return s.repo.FindByID(ctx, userID)
}

func (s *userService) ValidateToken(tokenString string) (map[string]interface{}, error) {
	return jwt.ValidateToken(tokenString, s.jwtSecret)
}
