package api

import (
	"auth-service/src/config"
	"auth-service/src/service"
	"encoding/json"
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

func (h *Handler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "corpo da requisição inválido"})
		return
	}

	user, err := h.service.Register(r.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		if err.Error() == "email already exists" {
			WriteJSON(w, http.StatusConflict, map[string]string{"error": "email já está em uso"})
			return
		}
		WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	response := map[string]interface{}{
		"id":        user.ID,
		"name":      user.Name,
		"email":     user.Email,
		"createdAt": user.CreatedAt,
	}
	WriteJSON(w, http.StatusCreated, response)
}

func (h *Handler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "corpo da requisição inválido"})
		return
	}

	token, err := h.service.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}
	WriteJSON(w, http.StatusOK, map[string]string{"token": token})
}

func (h *Handler) HandleGetProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userIDKey).(string)
	user, err := h.service.GetProfile(r.Context(), userID)
	if err != nil {
		WriteJSON(w, http.StatusNotFound, map[string]string{"error": "usuário não encontrado"})
		return
	}

	response := map[string]interface{}{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
	}
	WriteJSON(w, http.StatusOK, response)
}

func (h *Handler) HandleAuthValidate(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "corpo da requisição inválido"})
		return
	}

	claims, err := h.service.ValidateToken(req.Token)
	if err != nil {
		WriteJSON(w, http.StatusUnauthorized, map[string]bool{"valid": false})
		return
	}

	response := map[string]interface{}{
		"valid":  true,
		"userId": claims["sub"],
		"email":  claims["email"],
	}
	WriteJSON(w, http.StatusOK, response)
}

func WriteJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
