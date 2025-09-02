package api

import (
	"auth-service/src/config"
	"auth-service/src/domain"
	"auth-service/src/service"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandleRegister_Success(t *testing.T) {
	// Arrange
	mockService := new(service.UserServiceMock)
	handler := NewHandler(mockService, &config.Config{})

	// Prepara os dados de entrada
	requestBody := `{"name": "Test User", "email": "test@example.com", "password": "password123"}`
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(requestBody))
	rr := httptest.NewRecorder()

	// Programa o mock para retornar sucesso
	mockService.On("Register", mock.Anything, "Test User", "test@example.com", "password123").
		Return(&domain.User{ID: "123", Name: "Test User", Email: "test@example.com"}, nil)

	// Act
	handler.HandleRegister(rr, req)

	// Assert
	assert.Equal(t, http.StatusCreated, rr.Code)
	mockService.AssertExpectations(t)

	// Verifica o corpo da resposta
	var responseBody map[string]interface{}
	json.Unmarshal(rr.Body.Bytes(), &responseBody)
	assert.Equal(t, "123", responseBody["id"])
}

func TestHandleRegister_EmailAlreadyExists(t *testing.T) {
	// Arrange
	mockService := new(service.UserServiceMock)
	handler := NewHandler(mockService, &config.Config{})

	requestBody := `{"name": "Test User", "email": "exists@example.com", "password": "password123"}`
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(requestBody))
	rr := httptest.NewRecorder()

	// Programa o mock para retornar o erro de e-mail existente
	mockService.On("Register", mock.Anything, mock.Anything, "exists@example.com", mock.Anything).
		Return(nil, domain.ErrEmailAlreadyExists)

	// Act
	handler.HandleRegister(rr, req)

	// Assert
	assert.Equal(t, http.StatusConflict, rr.Code) // Verifica o status HTTP
	mockService.AssertExpectations(t)

	// Verifica se o corpo do erro está correto
	var errorResponse ErrorResponse
	json.Unmarshal(rr.Body.Bytes(), &errorResponse)
	assert.Equal(t, "EMAIL_ALREADY_EXISTS", errorResponse.Code)
}

func TestJWTAuthMiddleware_Success(t *testing.T) {
	// Arrange
	mockService := new(service.UserServiceMock)
	cfg := &config.Config{JWTSecret: "test-secret"}
	handler := NewHandler(mockService, cfg)

	// Cria um servidor de teste com o middleware e um handler final
	router := chi.NewRouter()
	router.With(handler.JWTAuthMiddleware).Get("/profile", handler.HandleGetProfile)
	testServer := httptest.NewServer(router)
	defer testServer.Close()

	// Programa os mocks
	mockService.On("ValidateToken", "valid-token").
		Return(map[string]interface{}{"sub": "user-123"}, nil)
	mockService.On("GetProfile", mock.Anything, "user-123").
		Return(&domain.User{ID: "user-123", Name: "Profile User"}, nil)

	// Act
	req, _ := http.NewRequest(http.MethodGet, testServer.URL+"/profile", nil)
	req.Header.Set("Authorization", "Bearer valid-token")
	res, err := http.DefaultClient.Do(req)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	mockService.AssertExpectations(t)
}
