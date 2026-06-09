package transport

import (
	"encoding/json"
	"net/http"

	"core-auth-org/internal/modules/org/service"
	"core-auth-org/internal/platform/server"

	"github.com/google/uuid"
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

	org, err := h.orgService.Create(r.Context(), req.Name)
	if err != nil {
		server.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	server.JSON(w, http.StatusCreated, org)
}

type createUnitRequest struct {
	OrganizationID uuid.UUID `json:"organization_id"`
	Name           string    `json:"name"`
}

func (h *OrgHandler) CreateUnit(w http.ResponseWriter, r *http.Request) {
	var req createUnitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		server.Error(w, http.StatusBadRequest, "payload inválido")
		return
	}

	unit, err := h.orgService.CreateUnit(r.Context(), req.OrganizationID, req.Name)
	if err != nil {
		server.Error(w, http.StatusInternalServerError, "falha ao criar filial")
		return
	}

	server.JSON(w, http.StatusCreated, unit)
}

func (h *OrgHandler) ListUnits(w http.ResponseWriter, r *http.Request) {
	orgIDStr := r.URL.Query().Get("organization_id")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		server.Error(w, http.StatusBadRequest, "organization_id inválido")
		return
	}

	units, err := h.orgService.ListUnits(r.Context(), orgID)
	if err != nil {
		server.Error(w, http.StatusInternalServerError, "falha ao buscar filiais")
		return
	}

	server.JSON(w, http.StatusOK, units)
}
