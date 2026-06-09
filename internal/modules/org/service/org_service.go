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

func (s *OrgService) CreateOrganization(ctx context.Context, name string, creatorID uuid.UUID) (*domain.Organization, error) {
	name = strings.TrimSpace(name)
	if len(name) < 3 {
		return nil, domain.ErrInvalidOrgName
	}

	// Gera um slug simples (ex: "Zanvexis Corp" -> "zanvexis-corp")
	slug := strings.ToLower(strings.ReplaceAll(name, " ", "-"))

	// 1. Cria a organização no banco
	dbOrg, err := s.repo.CreateOrganization(ctx, repository.CreateOrganizationParams{
		Name: name,
		Slug: slug,
	})
	if err != nil {
		return nil, domain.ErrOrgCreation
	}

	// 2. Adiciona o usuário criador como membro da organização
	err = s.repo.AddUserToOrganization(ctx, repository.AddUserToOrganizationParams{
		OrganizationID: dbOrg.ID,
		UserID:         creatorID,
	})
	if err != nil {
		return nil, err
	}

	return &domain.Organization{
		ID:        dbOrg.ID,
		Name:      dbOrg.Name,
		Slug:      dbOrg.Slug,
		CreatedAt: dbOrg.CreatedAt, // Correção: A propriedade já é nativa do tipo time.Time
		UpdatedAt: dbOrg.UpdatedAt, // Correção: A propriedade já é nativa do tipo time.Time
	}, nil
}