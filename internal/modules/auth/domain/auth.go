package domain

import "errors"

var (
	ErrInvalidCredentials = errors.New("email ou senha inválidos")
)