package service

import (
	"context"
	"net/mail"

	"core-auth-org/internal/modules/users/domain"
	"core-auth-org/internal/modules/users/repository"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo repository.Querier
}

func NewUserService(repo repository.Querier) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(ctx context.Context, email, password string) (*domain.User, error) {
	if _, err := mail.ParseAddress(email); err != nil {
		return nil, domain.ErrInvalidEmail
	}

	if len(password) < 8 {
		return nil, domain.ErrPasswordTooShort
	}

	_, err := s.repo.GetUserByEmail(ctx, email)
	if err == nil {
		return nil, domain.ErrEmailAlreadyInUse
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	dbUser, err := s.repo.CreateUser(ctx, repository.CreateUserParams{
		Email:        email,
		PasswordHash: string(hash),
	})
	if err != nil {
		return nil, err
	}

	return &domain.User{
		ID:           dbUser.ID,
		Email:        dbUser.Email,
		PasswordHash: dbUser.PasswordHash,
		IsActive:     dbUser.IsActive,
		CreatedAt:    dbUser.CreatedAt.Time,
		UpdatedAt:    dbUser.UpdatedAt.Time,
	}, nil
}