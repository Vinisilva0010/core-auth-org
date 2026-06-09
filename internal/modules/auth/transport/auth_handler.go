package transport

import (
	"encoding/json"
	"net/http"

	"core-auth-org/internal/modules/auth/service"
	"core-auth-org/internal/platform/server"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(svc *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: svc}
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		server.Error(w, http.StatusBadRequest, "payload inválido")
		return
	}

	accessToken, refreshToken, err := h.authService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		server.Error(w, http.StatusUnauthorized, "credenciais inválidas")
		return
	}

	server.JSON(w, http.StatusOK, map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req refreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		server.Error(w, http.StatusBadRequest, "payload inválido")
		return
	}

	accessToken, refreshToken, err := h.authService.Refresh(r.Context(), req.RefreshToken)
	if err != nil {
		server.Error(w, http.StatusUnauthorized, "sessão inválida ou expirada")
		return
	}

	server.JSON(w, http.StatusOK, map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

type logoutRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var req logoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		server.Error(w, http.StatusBadRequest, "payload inválido")
		return
	}

	_ = h.authService.Logout(r.Context(), req.RefreshToken)
	w.WriteHeader(http.StatusNoContent)
}

type forgotPasswordRequest struct {
	Email string `json:"email"`
}

func (h *AuthHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req forgotPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		server.Error(w, http.StatusBadRequest, "payload inválido")
		return
	}

	token, _ := h.authService.RequestPasswordReset(r.Context(), req.Email)
	
	// Retornando o token apenas em dev. Em prod, isso não deve ser enviado na resposta.
	server.JSON(w, http.StatusOK, map[string]string{
		"message": "Se o e-mail existir, um link de recuperação foi enviado.",
		"dev_token": token, 
	})
}

type resetPasswordRequest struct {
	Token       string `json:"token"`
	NewPassword string `json:"new_password"`
}

func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req resetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		server.Error(w, http.StatusBadRequest, "payload inválido")
		return
	}

	if err := h.authService.ResetPassword(r.Context(), req.Token, req.NewPassword); err != nil {
		server.Error(w, http.StatusBadRequest, "falha ao redefinir a senha")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
