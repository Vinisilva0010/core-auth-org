package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidRoleData = errors.New("dados do cargo inválidos")
	ErrRoleCreation    = errors.New("falha ao criar o cargo")
	ErrInvalidPermData = errors.New("dados de permissão inválidos")
)

type Role struct {
	ID             uuid.UUID `json:"id"`
	OrganizationID uuid.UUID `json:"organization_id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	CreatedAt      time.Time `json:"created_at"`
}

type Permission struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
