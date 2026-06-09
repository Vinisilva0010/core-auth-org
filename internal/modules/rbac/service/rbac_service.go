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

func (s *RBACService) AssignRole(ctx context.Context, orgID, userID, roleID uuid.UUID) error {
	return s.repo.AssignRoleToMember(ctx, repository.AssignRoleToMemberParams{
		RoleID:         roleID,
		OrganizationID: orgID,
		UserID:         userID,
	})
}

func (s *RBACService) CreatePermission(ctx context.Context, name string) (*domain.Permission, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, domain.ErrInvalidPermData
	}
	dbPerm, err := s.repo.CreatePermission(ctx, name)
	if err != nil {
		return nil, err
	}
	return &domain.Permission{
		ID:   dbPerm.ID,
		Name: dbPerm.Name,
	}, nil
}

func (s *RBACService) AssignPermissionToRole(ctx context.Context, roleID, permID uuid.UUID) error {
	return s.repo.AssignPermissionToRole(ctx, repository.AssignPermissionToRoleParams{
		RoleID:       roleID,
		PermissionID: permID,
	})
}
