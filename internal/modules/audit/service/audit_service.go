package service

import (
	"context"

	"core-auth-org/internal/modules/audit/domain"
	"core-auth-org/internal/modules/audit/repository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type AuditService struct {
	repo repository.Querier
}

func NewAuditService(repo repository.Querier) *AuditService {
	return &AuditService{repo: repo}
}

func (s *AuditService) LogEvent(ctx context.Context, orgID, userID uuid.UUID, action, resource string, details []byte, ip string) error {
	var pgOrgID, pgUserID pgtype.UUID
	
	if orgID != uuid.Nil {
		pgOrgID = pgtype.UUID{Bytes: orgID, Valid: true}
	}
	if userID != uuid.Nil {
		pgUserID = pgtype.UUID{Bytes: userID, Valid: true}
	}

	_, err := s.repo.CreateAuditLog(ctx, repository.CreateAuditLogParams{
		OrganizationID: pgOrgID,
		UserID:         pgUserID,
		Action:         action,
		Resource:       resource,
		Details:        details,
		IpAddress:      pgtype.Text{String: ip, Valid: ip != ""},
	})
	return err
}

func (s *AuditService) ListByOrg(ctx context.Context, orgID uuid.UUID) ([]domain.AuditLog, error) {
	dbLogs, err := s.repo.ListAuditLogsByOrg(ctx, pgtype.UUID{Bytes: orgID, Valid: true})
	if err != nil {
		return nil, err
	}

	var logs []domain.AuditLog
	for _, l := range dbLogs {
		logs = append(logs, domain.AuditLog{
			ID:             l.ID,
			OrganizationID: l.OrganizationID.Bytes,
			UserID:         l.UserID.Bytes,
			Action:         l.Action,
			Resource:       l.Resource,
			Details:        l.Details,
			IPAddress:      l.IpAddress.String,
			CreatedAt:      l.CreatedAt,
		})
	}
	return logs, nil
}
