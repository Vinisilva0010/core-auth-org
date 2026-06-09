package domain

import (
	"time"

	"github.com/google/uuid"
)

type AuditLog struct {
	ID             uuid.UUID `json:"id"`
	OrganizationID uuid.UUID `json:"organization_id,omitempty"`
	UserID         uuid.UUID `json:"user_id,omitempty"`
	Action         string    `json:"action"`
	Resource       string    `json:"resource"`
	Details        []byte    `json:"details,omitempty"`
	IPAddress      string    `json:"ip_address,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
}
