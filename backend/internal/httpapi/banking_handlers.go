package httpapi

import (
	"net/http"
	"time"

	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/openbanking"
)

func connectionJSON(c *openbanking.Connection) map[string]any {
	out := map[string]any{
		"id":         c.ID,
		"provider":   c.Provider,
		"status":     c.Status,
		"created_at": c.CreatedAt,
		"updated_at": c.UpdatedAt,
	}
	if c.ConsentExpiresAt != nil {
		out["consent_expires_at"] = c.ConsentExpiresAt.Format(time.RFC3339)
	}
	if c.LastSyncAt != nil {
		out["last_sync_at"] = c.LastSyncAt.Format(time.RFC3339)
	}
	return out
}

func handleBankingStatus(svc *openbanking.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if svc == nil {
			writeError(w, http.StatusServiceUnavailable, "open banking indisponível")
			return
		}
		userID, _ := r.Context().Value(userIDKey).(string)
		c, err := svc.Status(r.Context(), userID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "falha ao ler ligação")
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"connection": connectionJSON(c)})
	}
}

func handleBankingConnect(svc *openbanking.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if svc == nil {
			writeError(w, http.StatusServiceUnavailable, "open banking indisponível")
			return
		}
		userID, _ := r.Context().Value(userIDKey).(string)
		c, err := svc.Connect(r.Context(), userID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "falha ao ligar banco")
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"connection": connectionJSON(c)})
	}
}

func handleBankingSync(svc *openbanking.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if svc == nil {
			writeError(w, http.StatusServiceUnavailable, "open banking indisponível")
			return
		}
		userID, _ := r.Context().Value(userIDKey).(string)
		txs, conn, err := svc.Sync(r.Context(), userID)
		if err == openbanking.ErrNotConnected {
			writeError(w, http.StatusConflict, "liga o banco antes de sincronizar")
			return
		}
		if err != nil {
			writeError(w, http.StatusInternalServerError, "falha na sincronização")
			return
		}
		out := make([]map[string]any, 0, len(txs))
		for _, tx := range txs {
			out = append(out, map[string]any{
				"id":          tx.ID,
				"amount":      tx.Amount,
				"currency":    tx.Currency,
				"booked_at":   tx.BookedAt.Format(time.RFC3339),
				"description": tx.Description,
				"merchant_ref": tx.MerchantRef,
			})
		}
		writeJSON(w, http.StatusOK, map[string]any{
			"connection":   connectionJSON(conn),
			"transactions": out,
			"provider":     svc.Provider.Name(),
		})
	}
}
