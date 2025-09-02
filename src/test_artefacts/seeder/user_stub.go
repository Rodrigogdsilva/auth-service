package stubs

import (
	"auth-service/src/domain"
	"time"

	"github.com/google/uuid"
	"github.com/jaswdr/faker"
)

type UserStub struct {
	user *domain.User
}

func NewUserStub() *UserStub {
	f := faker.New()
	return &UserStub{
		user: &domain.User{
			ID:           uuid.NewString(),
			Name:         f.Person().Name(),
			Email:        f.Internet().Email(),
			PasswordHash: "hashed-password",
			CreatedAt:    time.Now(),
		},
	}
}

func (s *UserStub) WithEmail(email string) *UserStub {
	s.user.Email = email
	return s
}

func (s *UserStub) Get() *domain.User {
	return s.user
}
