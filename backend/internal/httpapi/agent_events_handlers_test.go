package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAgentEventsRoute_NotRegisteredWithoutDeps(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/agent/events", nil)
	rec := httptest.NewRecorder()
	NewRouter(Deps{}).ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status=%d", rec.Code)
	}
}
