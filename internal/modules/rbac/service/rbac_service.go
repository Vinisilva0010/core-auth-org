package service

import (
	"context"
	"strings"

	"core-auth-org/internal/modules/rbac/domain"
	"core-auth-org/internal/modules/rbac/repository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type RBACService struct {
	repo repository.Querier
}

func NewRBACService(repo repository.Querier) *RBACService {
	return &RBACService{repo: repo}
}

func (s *RBACService) CreateRole(ctx context.Context, orgID uuid.UUID, name, description string) (*domain.Role, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, domain.ErrInvalidRoleData
	}

	dbRole, err := s.repo.CreateRole(ctx, repository.CreateRoleParams{
		OrganizationID: orgID,
		Name:           name,
		Description:    pgtype.Text{String: description, Valid: description != ""},
	})
	if err != nil {
		return nil, domain.ErrRoleCreation
	}

	return &domain.Role{
		ID:             dbRole.ID,
		OrganizationID: dbRole.OrganizationID,
		Name:           dbRole.Name,
		Description:    dbRole.Description.String,
		CreatedAt:      dbRole.CreatedAt,
	}, nil
}
