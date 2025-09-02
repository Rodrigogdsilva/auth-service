package domain

import (
	"errors"
)

var (
	ErrUserNotFound          = errors.New("user not found")
	ErrEmailAlreadyExists    = errors.New("email already exists")
	ErrInvalidCredentials    = errors.New("invalid credentials")
	ErrParametersMissing     = errors.New("name, email and password are required")
	ErrPasswordTooShort      = errors.New("password must be at least 8 characters long")
	ErrCryptographyFailure   = errors.New("failed to encrypt password")
	ErrJwtSecretMissing      = errors.New("JWT secret is not configured")
	ErrInvalidToken          = errors.New("invalid authentication token")
	ErrUnexpected            = errors.New("an unexpected error occurred")
	ErrInvalidRequestBody    = errors.New("invalid request body")
	ErrFailedHashingPassword = errors.New("failed to hash password")
)
