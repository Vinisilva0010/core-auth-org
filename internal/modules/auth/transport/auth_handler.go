package transport

import (
	"encoding/json"
	"errors"
	"net/http"

	"core-auth-org/internal/modules/auth/domain"
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

type loginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		server.Error(w, http.StatusBadRequest, "payload inválido")
		return
	}

	ipAddress := r.RemoteAddr
	userAgent := r.UserAgent()

	accessToken, refreshToken, err := h.authService.Login(r.Context(), req.Email, req.Password, ipAddress, userAgent)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			server.Error(w, http.StatusUnauthorized, err.Error())
			return
		}
		server.Error(w, http.StatusInternalServerError, "erro interno durante o login")
		return
	}

	server.JSON(w, http.StatusOK, loginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

type tokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type refreshResponse struct {
	AccessToken string `json:"access_token"`
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req tokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		server.Error(w, http.StatusBadRequest, "payload inválido")
		return
	}

	accessToken, err := h.authService.RefreshToken(r.Context(), req.RefreshToken)
	if err != nil {
		server.Error(w, http.StatusUnauthorized, "refresh token inválido ou expirado")
		return
	}

	server.JSON(w, http.StatusOK, refreshResponse{AccessToken: accessToken})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var req tokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		server.Error(w, http.StatusBadRequest, "payload inválido")
		return
	}

	if err := h.authService.Logout(r.Context(), req.RefreshToken); err != nil {
		server.Error(w, http.StatusInternalServerError, "erro ao realizar logout")
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204: Sucesso sem corpo na resposta
}