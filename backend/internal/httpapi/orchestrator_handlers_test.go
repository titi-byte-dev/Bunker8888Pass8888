package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOrchestratorStatus_NotRegisteredWithoutDeps(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/agent/orchestrator/status", nil)
	rec := httptest.NewRecorder()
	NewRouter(Deps{}).ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status=%d", rec.Code)
	}
}
