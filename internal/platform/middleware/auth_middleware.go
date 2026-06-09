package middleware

import (
	"context"
	"net/http"
	"strings"

	"core-auth-org/internal/platform/server"
	"core-auth-org/internal/platform/token"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type contextKey string

const UserIDKey contextKey = "user_id"

// RequireAuth é o middleware que bloqueia requisições sem um JWT válido.
func RequireAuth(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				server.Error(w, http.StatusUnauthorized, "token de autenticação não fornecido")
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				server.Error(w, http.StatusUnauthorized, "formato de token inválido. Use: Bearer <token>")
				return
			}

			tokenString := parts[1]

			// Parse e validação da assinatura do token
			parsedToken, err := jwt.ParseWithClaims(tokenString, &token.Claims{}, func(t *jwt.Token) (interface{}, error) {
				return []byte(jwtSecret), nil
			})

			if err != nil || !parsedToken.Valid {
				server.Error(w, http.StatusUnauthorized, "token inválido ou expirado")
				return
			}

			claims, ok := parsedToken.Claims.(*token.Claims)
			if !ok {
				server.Error(w, http.StatusUnauthorized, "payload do token inválido")
				return
			}

			// Injeta o UserID no contexto da requisição para os próximos handlers usarem
			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserID extrai o UUID do usuário do contexto da requisição.
// Retorna erro se tentar extrair de uma rota que não passou pelo RequireAuth.
func GetUserID(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(UserIDKey).(uuid.UUID)
	return id, ok
}