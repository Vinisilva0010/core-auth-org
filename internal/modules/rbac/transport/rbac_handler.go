package transport

import (
	"encoding/json"
	"net/http"

	"core-auth-org/internal/modules/rbac/service"
	"core-auth-org/internal/platform/server"

	"github.com/google/uuid"
)

type RBACHandler struct {
	rbacService *service.RBACService
}

func NewRBACHandler(svc *service.RBACService) *RBACHandler {
	return &RBACHandler{rbacService: svc}
}

type createRoleRequest struct {
	OrganizationID uuid.UUID `json:"organization_id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
}

func (h *RBACHandler) CreateRole(w http.ResponseWriter, r *http.Request) {
	var req createRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		server.Error(w, http.StatusBadRequest, "payload inválido")
		return
	}

	role, err := h.rbacService.CreateRole(r.Context(), req.OrganizationID, req.Name, req.Description)
	if err != nil {
		server.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	server.JSON(w, http.StatusCreated, role)
}

type assignRoleRequest struct {
	OrganizationID uuid.UUID `json:"organization_id"`
	UserID         uuid.UUID `json:"user_id"`
	RoleID         uuid.UUID `json:"role_id"`
}

func (h *RBACHandler) AssignRole(w http.ResponseWriter, r *http.Request) {
	var req assignRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		server.Error(w, http.StatusBadRequest, "payload inválido")
		return
	}

	if err := h.rbacService.AssignRole(r.Context(), req.OrganizationID, req.UserID, req.RoleID); err != nil {
		server.Error(w, http.StatusInternalServerError, "falha ao atribuir cargo")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
