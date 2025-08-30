package jwt

import (
	"auth-service/src/domain"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateToken(user *domain.User, secret string) (string, error) {
	if secret == "" {
		return "", errors.New("segredo JWT não configurado")
	}
	claims := &jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ValidateToken(tokenString, secret string) (jwt.MapClaims, error) {
	if secret == "" {
		return nil, errors.New("segredo JWT não configurado")
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de assinatura inesperado: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("token inválido")
}
