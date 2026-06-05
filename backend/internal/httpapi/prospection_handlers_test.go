package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProspectionRoute_NotRegisteredWithoutDeps(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/agent/prospection/run", nil)
	rec := httptest.NewRecorder()
	NewRouter(Deps{}).ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status=%d", rec.Code)
	}
}
