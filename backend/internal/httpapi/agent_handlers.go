package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/agent"
)

type agentRunRequest struct {
	AgentID string          `json:"agent_id"`
	Input   json.RawMessage `json:"input"`
}

func handleListAgentTools(reg *agent.Registry) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]any{
			"tools": reg.ListDescriptors(),
		})
	}
}

func handleRunAgentTool(run *agent.Runner) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name := r.PathValue("name")
		if name == "" {
			writeError(w, http.StatusBadRequest, "nome da tool em falta")
			return
		}
		var req agentRunRequest
		if err := decodeJSON(r, &req); err != nil {
			writeError(w, http.StatusBadRequest, "JSON inválido")
			return
		}
		if req.AgentID == "" {
			req.AgentID = "manual"
		}
		userID, _ := r.Context().Value(userIDKey).(string)
		out, err := run.Run(r.Context(), name, agent.Request{
			UserID:  userID,
			AgentID: req.AgentID,
			Input:   req.Input,
		})
		if mapAgentError(w, err) {
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{
			"tool":   name,
			"output": json.RawMessage(out),
		})
	}
}

func mapAgentError(w http.ResponseWriter, err error) bool {
	if err == nil {
		return false
	}
	switch {
	case errors.Is(err, agent.ErrToolNotFound):
		writeError(w, http.StatusNotFound, "tool não encontrada")
	case errors.Is(err, agent.ErrInvalidToolInput),
		errors.Is(err, agent.ErrExternalAsCommand):
		writeError(w, http.StatusBadRequest, err.Error())
	case errors.Is(err, agent.ErrPermissionDenied):
		writeError(w, http.StatusForbidden, "permissão negada")
	default:
		writeError(w, http.StatusInternalServerError, "falha ao executar tool")
	}
	return true
}
