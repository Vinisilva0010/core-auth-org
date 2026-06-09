package transport

import (
	"encoding/json"
	"errors"
	"net/http"

	"core-auth-org/internal/modules/users/domain"
	"core-auth-org/internal/modules/users/service"
	"core-auth-org/internal/platform/server"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{userService: svc}
}

// DTOs privados do handler: o que entra e o que sai da API
type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type registerResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		server.Error(w, http.StatusBadRequest, "payload de requisição inválido")
		return
	}

	user, err := h.userService.CreateUser(r.Context(), req.Email, req.Password)
	if err != nil {
		// Mapeamento de erros de domínio para status HTTP
		switch {
		case errors.Is(err, domain.ErrEmailAlreadyInUse):
			server.Error(w, http.StatusConflict, err.Error())
		case errors.Is(err, domain.ErrInvalidEmail), errors.Is(err, domain.ErrPasswordTooShort):
			server.Error(w, http.StatusBadRequest, err.Error())
		default:
			server.Error(w, http.StatusInternalServerError, "erro interno ao criar usuário")
		}
		return
	}

	resp := registerResponse{
		ID:        user.ID.String(),
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	server.JSON(w, http.StatusCreated, resp)
}