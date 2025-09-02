package jwt

import (
	"auth-service/src/domain"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateAndValidateToken(t *testing.T) {
	// Arrange
	user := &domain.User{
		ID:    "a1b2c3d4-e5f6-4a7b-8c9d-0f1a2b3c4d5e",
		Email: "test@example.com",
	}
	secret := "my-super-secret-key-for-testing"

	// Act: Cria o token
	tokenString, err := CreateToken(user, secret)

	// Assert: Verifica se a criação foi bem-sucedida
	require.NoError(t, err)
	require.NotEmpty(t, tokenString)

	// Act: Valida o token recém-criado
	claims, err := ValidateToken(tokenString, secret)

	// Assert: Verifica se a validação foi bem-sucedida e se os dados estão corretos
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, user.ID, claims["sub"])
	assert.Equal(t, user.Email, claims["email"])

	// Verifica a data de expiração (exp)
	exp, ok := claims["exp"].(float64)
	assert.True(t, ok)
	assert.Greater(t, exp, float64(time.Now().Unix()))
}

func TestValidateToken_InvalidSignature(t *testing.T) {
	// Arrange
	user := &domain.User{ID: "user-id", Email: "test@example.com"}
	secret1 := "secret-one"
	secret2 := "secret-two" // secret diferente

	tokenString, err := CreateToken(user, secret1)
	require.NoError(t, err)

	// Act: Tenta validar o token com o secret errado
	claims, err := ValidateToken(tokenString, secret2)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, claims)
	assert.Contains(t, err.Error(), "signature is invalid")
}

func TestValidateToken_MissingSecret(t *testing.T) {
	// Act: Tenta validar com um secret vazio
	_, err := ValidateToken("any-token", "")

	// Assert
	assert.Error(t, err)
	assert.ErrorIs(t, err, domain.ErrJwtSecretMissing)
}
