package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/agent"
)

func TestAgentRoutes_NotRegisteredWithoutDeps(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/agent/tools", nil)
	rec := httptest.NewRecorder()
	NewRouter(Deps{}).ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status=%d, esperado 404", rec.Code)
	}
}

func TestHandleListAgentTools(t *testing.T) {
	reg := agent.NewDefaultRegistry()
	req := httptest.NewRequest(http.MethodGet, "/api/agent/tools", nil)
	rec := httptest.NewRecorder()
	handleListAgentTools(reg)(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
}

func TestHandleRunAgentTool(t *testing.T) {
	reg := agent.NewDefaultRegistry()
	run := agent.NewRunner(reg, agent.PermissivePolicy{})

	body, _ := json.Marshal(map[string]any{
		"agent_id": "test",
		"input":    map[string]string{"message": "ola"},
	})
	req := httptest.NewRequest(http.MethodPost, "/api/agent/tools/ping/run", bytes.NewReader(body))
	req.SetPathValue("name", "ping")
	req = req.WithContext(context.WithValue(req.Context(), userIDKey, "user-1"))
	rec := httptest.NewRecorder()
	handleRunAgentTool(run, nil, nil)(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
}
