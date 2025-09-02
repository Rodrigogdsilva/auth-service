package api

import (
	"auth-service/src/config"
	"auth-service/src/domain"
	"auth-service/src/service"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type Handler struct {
	service service.UserService
	cfg     *config.Config
}

func NewHandler(svc service.UserService, cfg *config.Config) *Handler {
	return &Handler{
		service: svc,
		cfg:     cfg,
	}
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (h *Handler) handleError(w http.ResponseWriter, err error) {

	log.Printf("ERRO: %v", err)

	if errors.Is(err, domain.ErrEmailAlreadyExists) {
		WriteJSON(w, http.StatusConflict, ErrorResponse{Code: "EMAIL_ALREADY_EXISTS", Message: domain.ErrEmailAlreadyExists.Error()})
		return
	}
	if errors.Is(err, domain.ErrInvalidCredentials) {
		WriteJSON(w, http.StatusUnauthorized, ErrorResponse{Code: "INVALID_CREDENTIALS", Message: domain.ErrInvalidCredentials.Error()})
		return
	}
	if errors.Is(err, domain.ErrUserNotFound) {
		WriteJSON(w, http.StatusUnauthorized, ErrorResponse{Code: "USER_NOT_FOUND", Message: domain.ErrUserNotFound.Error()})
		return
	}
	if errors.Is(err, domain.ErrParametersMissing) {
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{Code: "MISSING_PARAMETERS", Message: domain.ErrParametersMissing.Error()})
	}
	if errors.Is(err, domain.ErrPasswordTooShort) {
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{Code: "INVALID_INPUT", Message: domain.ErrPasswordTooShort.Error()})
		return
	}
	WriteJSON(w, http.StatusInternalServerError, ErrorResponse{Code: "INTERNAL_SERVER_ERROR", Message: domain.ErrUnexpected.Error()})
}

func (h *Handler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{Code: "INVALID_REQUEST_BODY", Message: domain.ErrInvalidRequestBody.Error()})
		return
	}

	user, err := h.service.Register(r.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		h.handleError(w, err)
		return
	}

	response := map[string]interface{}{"id": user.ID, "name": user.Name, "email": user.Email, "createdAt": user.CreatedAt}
	WriteJSON(w, http.StatusCreated, response)
}

func (h *Handler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{Code: "INVALID_REQUEST_BODY", Message: domain.ErrInvalidRequestBody.Error()})
		return
	}

	token, err := h.service.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		h.handleError(w, err)
		return
	}

	WriteJSON(w, http.StatusOK, map[string]string{"token": token})
}

func (h *Handler) HandleGetProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userIDKey).(string)
	user, err := h.service.GetProfile(r.Context(), userID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	response := map[string]interface{}{"id": user.ID, "name": user.Name, "email": user.Email}
	WriteJSON(w, http.StatusOK, response)
}

func (h *Handler) HandleAuthValidate(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteJSON(w, http.StatusBadRequest, ErrorResponse{Code: "INVALID_REQUEST_BODY", Message: domain.ErrInvalidRequestBody.Error()})
		return
	}

	claims, err := h.service.ValidateToken(req.Token)
	if err != nil {
		WriteJSON(w, http.StatusUnauthorized, map[string]bool{"valid": false})
		return
	}

	response := map[string]interface{}{"valid": true, "userId": claims["sub"], "email": claims["email"]}
	WriteJSON(w, http.StatusOK, response)
}

func WriteJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
