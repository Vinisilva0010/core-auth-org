package service

import (
	"context"
	"time"

	"core-auth-org/internal/modules/auth/domain"
	authRepo "core-auth-org/internal/modules/auth/repository"
	usersRepo "core-auth-org/internal/modules/users/repository"
	"core-auth-org/internal/platform/token"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	authRepo  authRepo.Querier
	userRepo  usersRepo.Querier
	jwtSecret string
}

func NewAuthService(ar authRepo.Querier, ur usersRepo.Querier, secret string) *AuthService {
	return &AuthService{
		authRepo:  ar,
		userRepo:  ur,
		jwtSecret: secret,
	}
}

func (s *AuthService) Login(ctx context.Context, email, password, ipAddress, userAgent string) (string, string, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", "", domain.ErrInvalidCredentials
	}

	if !user.IsActive {
		return "", "", domain.ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", "", domain.ErrInvalidCredentials
	}

	accessToken, err := token.Generate(user.ID, s.jwtSecret, 15*time.Minute)
	if err != nil {
		return "", "", err
	}

	refreshToken := uuid.NewString()

	_, err = s.authRepo.CreateSession(ctx, authRepo.CreateSessionParams{
		UserID:       user.ID,
		RefreshToken: refreshToken,
		UserAgent:    pgtype.Text{String: userAgent, Valid: userAgent != ""},
		IpAddress:    pgtype.Text{String: ipAddress, Valid: ipAddress != ""},
		ExpiresAt:    pgtype.Timestamptz{Time: time.Now().Add(7 * 24 * time.Hour), Valid: true},
	})
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// RefreshToken valida a sessão no banco e emite um novo Access Token de 15 minutos.
func (s *AuthService) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	session, err := s.authRepo.GetSessionByRefreshToken(ctx, refreshToken)
	if err != nil {
		return "", domain.ErrInvalidCredentials
	}

	if session.IsRevoked {
		return "", domain.ErrInvalidCredentials
	}

	// Lendo a propriedade .Time de dentro do pgtype para validar a expiração
	if time.Now().After(session.ExpiresAt.Time) {
		return "", domain.ErrInvalidCredentials
	}

	accessToken, err := token.Generate(session.UserID, s.jwtSecret, 15*time.Minute)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

// Logout invalida o Refresh Token no banco, matando a sessão.
func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	return s.authRepo.RevokeSession(ctx, refreshToken)
}