// Command server é o ponto de entrada da API HTTP do AegisPass.
package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/auth"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/config"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/db"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/httpapi"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/realtime"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/sessions"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/users"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/vault"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	cfg := config.Load()

	// Construímos as dependências do router. Sem DATABASE_URL, arrancamos em
	// modo mínimo (só /healthz) — útil para smoke tests sem base de dados.
	var deps httpapi.Deps
	var pool *pgxpool.Pool

	if cfg.DatabaseURL != "" {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var err error
		pool, err = db.Connect(ctx, cfg.DatabaseURL)
		cancel()
		if err != nil {
			logger.Error("ligação à BD falhou", "err", err)
			os.Exit(1)
		}
		defer pool.Close()

		migCtx, migCancel := context.WithTimeout(context.Background(), 30*time.Second)
		if err := db.Migrate(migCtx, pool); err != nil {
			migCancel()
			logger.Error("migrações falharam", "err", err)
			os.Exit(1)
		}
		migCancel()

		userRepo := users.NewRepo(pool)
		sessionRepo := sessions.NewRepo(pool)
		const sessionTTL = 24 * 60 * 60
		hub := realtime.NewHub()
		deps = httpapi.Deps{
			Auth:  auth.NewService(userRepo, sessionRepo, sessionTTL),
			Vault: vault.NewRepo(pool),
			Hub:   hub,
		}
		logger.Info("base de dados ligada e migrada")
	} else {
		logger.Warn("AEGIS_DATABASE_URL vazio: a arrancar em modo mínimo (só /healthz)")
	}

	srv := &http.Server{
		Addr:    cfg.HTTPAddr,
		Handler: httpapi.NewRouter(deps),
	}

	go func() {
		logger.Info("servidor a arrancar", "addr", cfg.HTTPAddr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("falha no servidor HTTP", "err", err)
			os.Exit(1)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()
	logger.Info("sinal de paragem recebido; a encerrar...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("encerramento forçado", "err", err)
		os.Exit(1)
	}
	logger.Info("servidor encerrado em segurança")
}
