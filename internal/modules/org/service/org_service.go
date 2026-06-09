package service

import (
	"context"
	"strings"

	"core-auth-org/internal/modules/org/domain"
	"core-auth-org/internal/modules/org/repository"

	"github.com/google/uuid"
)

type OrgService struct {
	repo repository.Querier
}

func NewOrgService(repo repository.Querier) *OrgService {
	return &OrgService{repo: repo}
}

func (s *OrgService) Create(ctx context.Context, name string) (*domain.Organization, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, domain.ErrInvalidOrgData
	}

	dbOrg, err := s.repo.CreateOrganization(ctx, repository.CreateOrganizationParams{
		Name: name,
	})
	if err != nil {
		return nil, domain.ErrOrgCreation
	}

	return &domain.Organization{
		ID:        dbOrg.ID,
		Name:      dbOrg.Name,
		CreatedAt: dbOrg.CreatedAt,
	}, nil
}

func (s *OrgService) CreateUnit(ctx context.Context, orgID uuid.UUID, name string) (*domain.OrganizationUnit, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, domain.ErrInvalidOrgData
	}

	dbUnit, err := s.repo.CreateOrganizationUnit(ctx, repository.CreateOrganizationUnitParams{
		OrganizationID: orgID,
		Name:           name,
	})
	if err != nil {
		return nil, err
	}

	return &domain.OrganizationUnit{
		ID:             dbUnit.ID,
		OrganizationID: dbUnit.OrganizationID,
		Name:           dbUnit.Name,
		CreatedAt:      dbUnit.CreatedAt,
	}, nil
}

func (s *OrgService) ListUnits(ctx context.Context, orgID uuid.UUID) ([]domain.OrganizationUnit, error) {
	dbUnits, err := s.repo.ListOrganizationUnits(ctx, orgID)
	if err != nil {
		return nil, err
	}

	var units []domain.OrganizationUnit
	for _, u := range dbUnits {
		units = append(units, domain.OrganizationUnit{
			ID:             u.ID,
			OrganizationID: u.OrganizationID,
			Name:           u.Name,
			CreatedAt:      u.CreatedAt,
		})
	}
	return units, nil
}
