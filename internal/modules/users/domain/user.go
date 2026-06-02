package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrUserNotFound      = errors.New("usuário não encontrado")
	ErrEmailAlreadyInUse = errors.New("este email já está em uso")
	ErrInvalidEmail      = errors.New("formato de email inválido")
	ErrPasswordTooShort  = errors.New("a senha deve ter no mínimo 8 caracteres")
)

// User representa a entidade central de negócio. 
// Note que não tem tags JSON ou DB aqui. O domínio é puro.
type User struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
	IsActive     bool
	CreatedAt    time.Time
	UpdatedAt    time.Time
}