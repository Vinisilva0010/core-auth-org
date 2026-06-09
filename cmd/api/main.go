package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"core-auth-org/internal/platform/config"
	"core-auth-org/internal/platform/logger"
	customMiddleware "core-auth-org/internal/platform/middleware"
	"core-auth-org/internal/platform/server"

	auditRepo "core-auth-org/internal/modules/audit/repository"
	auditSvc "core-auth-org/internal/modules/audit/service"
	auditTransport "core-auth-org/internal/modules/audit/transport"

	authRepo "core-auth-org/internal/modules/auth/repository"
	authSvc "core-auth-org/internal/modules/auth/service"
	authTransport "core-auth-org/internal/modules/auth/transport"

	orgRepo "core-auth-org/internal/modules/org/repository"
	orgSvc "core-auth-org/internal/modules/org/service"
	orgTransport "core-auth-org/internal/modules/org/transport"

	rbacRepo "core-auth-org/internal/modules/rbac/repository"
	rbacSvc "core-auth-org/internal/modules/rbac/service"
	rbacTransport "core-auth-org/internal/modules/rbac/transport"

	usersRepo "core-auth-org/internal/modules/users/repository"
	usersSvc "core-auth-org/internal/modules/users/service"
	usersTransport "core-auth-org/internal/modules/users/transport"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg := config.Load()
	logger.Init(cfg.Env)

	ctx := context.Background()

	dbPool, err := pgxpool.New(ctx, cfg.DBURL)
	if err != nil {
		slog.Error("falha ao criar pool de conexões", "erro", err)
		os.Exit(1)
	}
	defer dbPool.Close()

	if err := dbPool.Ping(ctx); err != nil {
		slog.Error("falha ao conectar no banco de dados", "erro", err)
		os.Exit(1)
	}
	slog.Info("banco de dados conectado com sucesso")

	userRepoInstance := usersRepo.New(dbPool)
	authRepoInstance := authRepo.New(dbPool)
	orgRepoInstance := orgRepo.New(dbPool)
	rbacRepoInstance := rbacRepo.New(dbPool)
	auditRepoInstance := auditRepo.New(dbPool)

	userSvcInstance := usersSvc.NewUserService(userRepoInstance)
	authSvcInstance := authSvc.NewAuthService(authRepoInstance, userRepoInstance, cfg.JWTSecret)
	orgSvcInstance := orgSvc.NewOrgService(orgRepoInstance)
	rbacSvcInstance := rbacSvc.NewRBACService(rbacRepoInstance)
	auditSvcInstance := auditSvc.NewAuditService(auditRepoInstance)

	userHandler := usersTransport.NewUserHandler(userSvcInstance)
	authHandler := authTransport.NewAuthHandler(authSvcInstance)
	orgHandler := orgTransport.NewOrgHandler(orgSvcInstance)
	rbacHandler := rbacTransport.NewRBACHandler(rbacSvcInstance)
	auditHandler := auditTransport.NewAuditHandler(auditSvcInstance)

	r := chi.NewRouter()

	r.Use(chiMiddleware.RequestID)
	r.Use(chiMiddleware.RealIP)
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer)
	r.Use(chiMiddleware.Timeout(60 * time.Second))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		server.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/users", userHandler.Register)
		r.Post("/auth/login", authHandler.Login)
		r.Post("/auth/refresh", authHandler.Refresh)
		r.Post("/auth/logout", authHandler.Logout)

		r.Group(func(r chi.Router) {
			r.Use(customMiddleware.RequireAuth(cfg.JWTSecret))

			r.Get("/users/me", func(w http.ResponseWriter, r *http.Request) {
				userID, _ := customMiddleware.GetUserID(r.Context())
				server.JSON(w, http.StatusOK, map[string]string{
					"status":  "autenticado",
					"user_id": userID.String(),
				})
			})

			r.Post("/organizations", orgHandler.Create)
			
			r.Post("/roles", rbacHandler.CreateRole)
			r.Post("/roles/assign", rbacHandler.AssignRole)

			r.Get("/audit", auditHandler.ListLogs)
		})
	})

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		slog.Info("servidor rodando", "porta", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("erro no servidor", "erro", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	slog.Info("desligando servidor...")

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctxShutdown); err != nil {
		slog.Error("erro no encerramento forçado", "erro", err)
	}
	slog.Info("servidor encerrado com segurança")
}
