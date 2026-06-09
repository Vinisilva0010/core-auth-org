package transport

import (
	"net/http"

	"core-auth-org/internal/modules/audit/service"
	"core-auth-org/internal/platform/server"

	"github.com/google/uuid"
)

type AuditHandler struct {
	auditService *service.AuditService
}

func NewAuditHandler(svc *service.AuditService) *AuditHandler {
	return &AuditHandler{auditService: svc}
}

func (h *AuditHandler) ListLogs(w http.ResponseWriter, r *http.Request) {
	orgIDStr := r.URL.Query().Get("organization_id")
	orgID, err := uuid.Parse(orgIDStr)
	if err != nil {
		server.Error(w, http.StatusBadRequest, "organization_id inválido")
		return
	}

	logs, err := h.auditService.ListByOrg(r.Context(), orgID)
	if err != nil {
		server.Error(w, http.StatusInternalServerError, "falha ao buscar logs")
		return
	}

	server.JSON(w, http.StatusOK, logs)
}
