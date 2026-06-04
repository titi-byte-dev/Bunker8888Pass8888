package httpapi

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/auth"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/realtime"
)

// wsUpgrader faz o "upgrade" HTTP → WebSocket.
//
// ⚠️ Segurança: CheckOrigin valida o header Origin. Em produção deve ser
// restrito aos domínios do frontend; por agora aceitamos qualquer origem
// (desenvolvimento local). Refinar em INFRA-001 com lista de origens.
var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool { return true },
}

// handleVaultWS mantém uma ligação WebSocket para sync do cofre.
//
// Autenticação: o browser não permite headers customizados no WebSocket nativo,
// por isso o token vai na query string (?token=...). ⚠️ Usar sempre WSS em
// produção para o token não viajar em claro.
func handleVaultWS(authSvc *auth.Service, hub *realtime.Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")
		if token == "" {
			writeError(w, http.StatusUnauthorized, "token em falta")
			return
		}
		userID, err := authSvc.Authenticate(r.Context(), token)
		if err != nil {
			writeError(w, http.StatusUnauthorized, "sessão inválida")
			return
		}

		conn, err := wsUpgrader.Upgrade(w, r, nil)
		if err != nil {
			slog.Debug("upgrade WebSocket falhou", "err", err)
			return
		}

		client := realtime.NewClient()
		hub.Register(userID, client)
		defer func() {
			hub.Unregister(userID, client)
			client.Close()
			_ = conn.Close()
		}()

		// Goroutine de escrita: lê da fila e envia frames WebSocket.
		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()

		go writePump(ctx, conn, client)

		// Goroutine de leitura: responde a pings e detecta desligamento.
		readPump(ctx, conn, cancel)
	}
}

func writePump(ctx context.Context, conn *websocket.Conn, c *realtime.Client) {
	ticker := time.NewTicker(realtime.PingPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-c.Send():
			_ = conn.SetWriteDeadline(time.Now().Add(realtime.WriteWait))
			if !ok {
				_ = conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}
		case <-ticker.C:
			_ = conn.SetWriteDeadline(time.Now().Add(realtime.WriteWait))
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func readPump(ctx context.Context, conn *websocket.Conn, cancel context.CancelFunc) {
	defer cancel()
	_ = conn.SetReadDeadline(time.Now().Add(realtime.PongWait))
	conn.SetPongHandler(func(string) error {
		_ = conn.SetReadDeadline(time.Now().Add(realtime.PongWait))
		return nil
	})
	for {
		if ctx.Err() != nil {
			return
		}
		if _, _, err := conn.ReadMessage(); err != nil {
			return
		}
	}
}
