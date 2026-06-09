package transport

import (
	"encoding/json"
	"net/http"

	"core-auth-org/internal/modules/org/service"
	"core-auth-org/internal/platform/middleware"
	"core-auth-org/internal/platform/server"
)

type OrgHandler struct {
	orgService *service.OrgService
}

func NewOrgHandler(svc *service.OrgService) *OrgHandler {
	return &OrgHandler{orgService: svc}
}

type createOrgRequest struct {
	Name string `json:"name"`
}

func (h *OrgHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createOrgRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		server.Error(w, http.StatusBadRequest, "payload inválido")
		return
	}

	// Extrai o ID do usuário injetado pelo middleware de autenticação
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		server.Error(w, http.StatusUnauthorized, "usuário não autenticado")
		return
	}

	org, err := h.orgService.CreateOrganization(r.Context(), req.Name, userID)
	if err != nil {
		server.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	server.JSON(w, http.StatusCreated, org)
}