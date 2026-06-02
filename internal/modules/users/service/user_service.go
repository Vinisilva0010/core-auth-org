package service

import (
	"context"
	"net/mail"

	"core-auth-org/internal/modules/users/domain"
	"core-auth-org/internal/modules/users/repository"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo repository.Querier // Usamos a interface gerada pelo sqlc para facilitar mocks no futuro
}

func NewUserService(repo repository.Querier) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(ctx context.Context, email, password string) (*domain.User, error) {
	// 1. Validação de domínio
	if _, err := mail.ParseAddress(email); err != nil {
		return nil, domain.ErrInvalidEmail
	}

	if len(password) < 8 {
		return nil, domain.ErrPasswordTooShort
	}

	// 2. Regra de negócio: evitar duplicidade
	// Verificamos se já existe. O repo retorna erro se não achar, então err == nil significa que achou.
	_, err := s.repo.GetUserByEmail(ctx, email)
	if err == nil {
		return nil, domain.ErrEmailAlreadyInUse
	}

	// 3. Segurança: Hash da senha
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 4. Persistência
	dbUser, err := s.repo.CreateUser(ctx, repository.CreateUserParams{
		Email:        email,
		PasswordHash: string(hash),
	})
	if err != nil {
		return nil, err
	}

	// 5. Retorna a entidade de domínio
	return &domain.User{
		ID:           dbUser.ID,
		Email:        dbUser.Email,
		PasswordHash: dbUser.PasswordHash,
		IsActive:     dbUser.IsActive,
		CreatedAt:    dbUser.CreatedAt,
		UpdatedAt:    dbUser.UpdatedAt,
	}, nil
}