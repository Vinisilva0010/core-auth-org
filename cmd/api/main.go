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
	"core-auth-org/internal/platform/database"
	"core-auth-org/internal/platform/logger"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg := config.Load()
	logger.Init(cfg.Env)

	// Contexto inicial para conexão com o banco
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbPool, err := database.Connect(ctx, cfg.DBURL)
	if err != nil {
		slog.Error("falha ao conectar no banco de dados", "erro", err)
		os.Exit(1)
	}
	defer dbPool.Close()

	slog.Info("banco de dados conectado com sucesso")

	r := chi.NewRouter()
	
	// Middlewares globais nativos do chi
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// Healthcheck
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	// Inicia o servidor em uma goroutine
	go func() {
		slog.Info("servidor rodando", "porta", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("erro no servidor", "erro", err)
			os.Exit(1)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("sinal de parada recebido, encerrando o servidor...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("erro ao encerrar servidor", "erro", err)
		os.Exit(1)
	}

	slog.Info("servidor encerrado com segurança")
}