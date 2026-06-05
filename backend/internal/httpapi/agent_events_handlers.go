package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/eventbus"
)

func handleListAgentEvents(store *eventbus.PGStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if store == nil {
			writeJSON(w, http.StatusOK, map[string]any{"events": []any{}})
			return
		}
		userID, _ := r.Context().Value(userIDKey).(string)
		records, err := store.ListRecent(r.Context(), userID, 50)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "falha ao listar eventos")
			return
		}
		out := make([]map[string]any, 0, len(records))
		for _, rec := range records {
			out = append(out, map[string]any{
				"id":         rec.ID,
				"type":       rec.Type,
				"source":     rec.Source,
				"payload":    json.RawMessage(rec.Payload),
				"created_at": rec.CreatedAt,
			})
		}
		writeJSON(w, http.StatusOK, map[string]any{"events": out})
	}
}
