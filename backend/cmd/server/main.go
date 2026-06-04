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

	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/config"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/httpapi"
)

func main() {
	// slog é o logger estruturado da biblioteca padrão (Go 1.21+). Logs
	// estruturados (chave=valor) são mais fáceis de pesquisar do que texto livre.
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	cfg := config.Load()

	srv := &http.Server{
		Addr:    cfg.HTTPAddr,
		Handler: httpapi.NewRouter(),
	}

	// Arrancamos o servidor numa goroutine para não bloquear a main, que vai
	// ficar à espera de um sinal de paragem (Ctrl+C / SIGTERM do Docker).
	go func() {
		logger.Info("servidor a arrancar", "addr", cfg.HTTPAddr)
		// ListenAndServe bloqueia até o servidor fechar. Quando fechamos de
		// propósito, devolve ErrServerClosed — que NÃO é um erro real.
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("falha no servidor HTTP", "err", err)
			os.Exit(1)
		}
	}()

	// signal.NotifyContext devolve um context que é cancelado quando chega um
	// dos sinais indicados. É a forma moderna de fazer graceful shutdown.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done() // bloqueia aqui até chegar o sinal
	logger.Info("sinal de paragem recebido; a encerrar...")

	// Damos um prazo para os pedidos em curso terminarem antes de matar tudo.
	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("encerramento forçado", "err", err)
		os.Exit(1)
	}
	logger.Info("servidor encerrado em segurança")
}
