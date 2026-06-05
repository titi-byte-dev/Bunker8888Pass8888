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
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/clidevices"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/config"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/db"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/httpapi"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/emergency"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/geofence"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/passkeys"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/recovery"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/realtime"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/security"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/sentinel"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/sessions"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/sharekeys"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/shifts"
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
		wipeSvc := security.NewWipeService(sessionRepo, hub, pool)
		shiftRepo := shifts.NewRepo(pool)
		geoRepo := geofence.NewRepo(pool)
		recoveryRepo := recovery.NewRepo(pool)
		emergencyRepo := emergency.NewRepo(pool)
		deviceRepo := clidevices.NewRepo(pool)
		passkeyRepo := passkeys.NewRepo(pool)
		sentinelRepo := sentinel.NewRepo(pool)
		shareKeysRepo := sharekeys.NewRepo(pool)

		var mtlsMat *config.MTLSMaterial
		if cfg.MTLSAutoDev || cfg.MTLSCACert != "" {
			var err error
			mtlsMat, err = config.LoadMTLSMaterial(cfg)
			if err != nil {
				logger.Error("mTLS material inválido", "err", err)
				os.Exit(1)
			}
		}

		var cliCA *clidevices.CA
		if mtlsMat != nil {
			cliCA = mtlsMat.CA
		}

		var passkeySvc *passkeys.Service
		if len(cfg.WebAuthnRPOrigins) > 0 && cfg.WebAuthnRPID != "" {
			var err error
			passkeySvc, err = passkeys.NewService(passkeys.Config{
				RPDisplayName: cfg.WebAuthnRPDisplayName,
				RPID:          cfg.WebAuthnRPID,
				RPOrigins:     cfg.WebAuthnRPOrigins,
			}, passkeyRepo, userRepo)
			if err != nil {
				logger.Error("WebAuthn inválido", "err", err)
				os.Exit(1)
			}
		}

		deps = httpapi.Deps{
			Auth:     auth.NewService(userRepo, sessionRepo, sessionTTL),
			Vault:    vault.NewRepo(pool),
			Hub:      hub,
			Wipe:     wipeSvc,
			Users:    userRepo,
			Shifts:   shiftRepo,
			Geofence: geoRepo,
			Recovery: recoveryRepo,
			Emergency: emergency.NewService(emergencyRepo, userRepo),
			Devices:  deviceRepo,
			CLIca:    cliCA,
			Passkeys:  passkeySvc,
			Sentinel:  sentinel.NewService(sentinelRepo),
			ShareKeys: shareKeysRepo,
			AdminKey:  cfg.AdminKey,
			Pool:     pool,
		}
		logger.Info("base de dados ligada e migrada")
	} else {
		logger.Warn("AEGIS_DATABASE_URL vazio: a arrancar em modo mínimo (só /healthz)")
	}

	srv := &http.Server{
		Addr:    cfg.HTTPAddr,
		Handler: httpapi.NewRouter(deps),
	}

	var mtlsSrv *http.Server
	if mtlsMat := loadMTLSForServer(cfg, deps, logger); mtlsMat != nil {
		mtlsSrv = mtlsMat
	}

	go func() {
		logger.Info("servidor a arrancar", "addr", cfg.HTTPAddr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("falha no servidor HTTP", "err", err)
			os.Exit(1)
		}
	}()

	if mtlsSrv != nil {
		go func() {
			logger.Info("servidor mTLS CLI a arrancar", "addr", mtlsSrv.Addr)
			if err := mtlsSrv.ListenAndServeTLS("", ""); err != nil && !errors.Is(err, http.ErrServerClosed) {
				logger.Error("falha no servidor mTLS", "err", err)
				os.Exit(1)
			}
		}()
	}

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
	if mtlsSrv != nil {
		if err := mtlsSrv.Shutdown(shutdownCtx); err != nil {
			logger.Error("encerramento mTLS forçado", "err", err)
			os.Exit(1)
		}
	}
	logger.Info("servidor encerrado em segurança")
}

func loadMTLSForServer(cfg config.Config, deps httpapi.Deps, logger *slog.Logger) *http.Server {
	if deps.Devices == nil || deps.Vault == nil {
		return nil
	}
	mtlsMat, err := config.LoadMTLSMaterial(cfg)
	if err != nil || mtlsMat == nil {
		if err != nil {
			logger.Warn("mTLS desactivado", "err", err)
		}
		return nil
	}
	if cfg.MTLSAddr == "" {
		return nil
	}
	tlsCfg, err := httpapi.MTLSTLSConfig(mtlsMat.ServerCertPEM, mtlsMat.ServerKeyPEM, mtlsMat.CACertPEM)
	if err != nil {
		logger.Warn("mTLS TLS config inválida", "err", err)
		return nil
	}
	return &http.Server{
		Addr:      cfg.MTLSAddr,
		Handler:   httpapi.NewMTLSRouter(deps),
		TLSConfig: tlsCfg,
	}
}
