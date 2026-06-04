// Package httpapi monta o router HTTP da API do AegisPass.
//
// Didático: separamos a construção do router (aqui) do arranque do processo
// (cmd/server). Assim o router pode ser testado isoladamente, sem abrir portas.
package httpapi

import (
	"encoding/json"
	"net/http"
	"time"
)

// NewRouter devolve o http.Handler com todas as rotas registadas.
//
// Em Go, http.Handler é uma interface com um único método ServeHTTP. O
// http.ServeMux (multiplexer) implementa-a e encaminha pedidos por padrão de URL.
func NewRouter() http.Handler {
	mux := http.NewServeMux()

	// Endpoint de "saúde" (health check): usado por orquestradores/Docker para
	// saber se o serviço está vivo. Não revela nada sensível.
	mux.HandleFunc("GET /healthz", handleHealth)

	return mux
}

// handleHealth responde com um JSON simples de estado.
func handleHealth(w http.ResponseWriter, _ *http.Request) {
	// Definir o Content-Type ANTES de escrever o corpo é importante: depois de
	// começar a escrever, os headers já foram enviados e não mudam.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(map[string]any{
		"status": "ok",
		"time":   time.Now().UTC().Format(time.RFC3339),
	})
}
