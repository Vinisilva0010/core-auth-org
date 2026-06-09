package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"time"

	"core-auth-org/internal/modules/auth/repository"
	usersRepo "core-auth-org/internal/modules/users/repository"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("credenciais inválidas")
	ErrSessionInvalid     = errors.New("sessão inválida ou expirada")
	ErrInvalidToken       = errors.New("token inválido ou expirado")
)

type AuthService struct {
	authRepo  repository.Querier
	userRepo  usersRepo.Querier
	jwtSecret string
}

func NewAuthService(ar repository.Querier, ur usersRepo.Querier, secret string) *AuthService {
	return &AuthService{authRepo: ar, userRepo: ur, jwtSecret: secret}
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, string, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", "", ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", "", ErrInvalidCredentials
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID.String(),
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
	})
	accessToken, _ := token.SignedString([]byte(s.jwtSecret))

	b := make([]byte, 32)
	rand.Read(b)
	refreshToken := base64.URLEncoding.EncodeToString(b)
	refreshExp := time.Now().Add(7 * 24 * time.Hour)

	// SESSÃO: Usa pgtype.Timestamptz exigido pelo compilador
	_, err = s.authRepo.CreateSession(ctx, repository.CreateSessionParams{
		UserID:       user.ID,
		RefreshToken: refreshToken,
		ExpiresAt:    pgtype.Timestamptz{Time: refreshExp, Valid: true},
	})
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (string, string, error) {
	session, err := s.authRepo.GetSessionByToken(ctx, refreshToken)
	if err != nil || session.IsRevoked || session.ExpiresAt.Time.Before(time.Now()) {
		return "", "", ErrSessionInvalid
	}

	_ = s.authRepo.RevokeSession(ctx, session.ID)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": session.UserID.String(),
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
	})
	newAccessToken, _ := token.SignedString([]byte(s.jwtSecret))

	b := make([]byte, 32)
	rand.Read(b)
	newRefreshToken := base64.URLEncoding.EncodeToString(b)
	refreshExp := time.Now().Add(7 * 24 * time.Hour)

	// SESSÃO: Usa pgtype.Timestamptz exigido pelo compilador
	_, err = s.authRepo.CreateSession(ctx, repository.CreateSessionParams{
		UserID:       session.UserID,
		RefreshToken: newRefreshToken,
		ExpiresAt:    pgtype.Timestamptz{Time: refreshExp, Valid: true},
	})

	return newAccessToken, newRefreshToken, err
}

func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	session, err := s.authRepo.GetSessionByToken(ctx, refreshToken)
	if err != nil {
		return nil
	}
	return s.authRepo.RevokeSession(ctx, session.ID)
}

func (s *AuthService) RequestPasswordReset(ctx context.Context, email string) (string, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", nil 
	}

	b := make([]byte, 32)
	rand.Read(b)
	resetToken := hex.EncodeToString(b)
	exp := time.Now().Add(1 * time.Hour)

	// RECUPERAÇÃO DE SENHA: Usa time.Time nativo exigido pelo compilador
	_, err = s.authRepo.CreatePasswordReset(ctx, repository.CreatePasswordResetParams{
		UserID:    user.ID,
		Token:     resetToken,
		ExpiresAt: exp,
	})
	if err != nil {
		return "", err
	}

	return resetToken, nil
}

func (s *AuthService) ResetPassword(ctx context.Context, token, newPassword string) error {
	resetReq, err := s.authRepo.GetPasswordResetByToken(ctx, token)
	// RECUPERAÇÃO DE SENHA: Usa time.Time nativo exigido pelo compilador
	if err != nil || resetReq.Used || resetReq.ExpiresAt.Before(time.Now()) {
		return ErrInvalidToken
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	err = s.userRepo.UpdateUserPassword(ctx, usersRepo.UpdateUserPasswordParams{
		ID:           resetReq.UserID,
		PasswordHash: string(hash),
	})
	if err != nil {
		return err
	}

	_ = s.authRepo.MarkPasswordResetUsed(ctx, resetReq.ID)
	_ = s.authRepo.RevokeAllUserSessions(ctx, resetReq.UserID)

	return nil
}
