package httpapi

import (
	"context"
	"errors"
	"net/http"

	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/eventbus"
)

func handleDecideOrchestratorAction(store *eventbus.PGStore, bus *eventbus.Bus, decision eventbus.Decision) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if store == nil {
			writeError(w, http.StatusServiceUnavailable, "aprovações indisponíveis")
			return
		}
		userID, _ := r.Context().Value(userIDKey).(string)
		suggestionID := r.PathValue("id")
		action, payload, err := eventbus.Decide(r.Context(), store, bus, userID, suggestionID, decision)
		if errors.Is(err, eventbus.ErrSuggestionNotFound) {
			writeError(w, http.StatusNotFound, "sugestão não encontrada")
			return
		}
		if errors.Is(err, eventbus.ErrAlreadyDecided) {
			writeError(w, http.StatusConflict, "sugestão já decidida")
			return
		}
		if err != nil {
			writeError(w, http.StatusInternalServerError, "falha ao registar decisão")
			return
		}
		status := "approved"
		if decision == eventbus.DecisionReject {
			status = "rejected"
		}
		writeJSON(w, http.StatusOK, map[string]any{
			"status":     status,
			"action":     action,
			"suggestion": suggestionID,
			"payload":    payload,
		})
	}
}

func enrichEventsWithApprovalStatus(ctx context.Context, store *eventbus.PGStore, userID string, out []map[string]any) error {
	if store == nil {
		return nil
	}
	decisions, err := store.DecisionMap(ctx, userID)
	if err != nil {
		return err
	}
	for _, row := range out {
		if row["type"] != eventbus.OrchestratorActionSuggested {
			continue
		}
		id, _ := row["id"].(string)
		if st := decisions[id]; st != "" {
			row["approval_status"] = st
		} else {
			row["approval_status"] = "pending"
		}
	}
	return nil
}
