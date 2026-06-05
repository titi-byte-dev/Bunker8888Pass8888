package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestApprovalRoutes_NotRegisteredWithoutDeps(t *testing.T) {
	for _, path := range []string{
		"/api/agent/orchestrator/actions/x/approve",
		"/api/agent/orchestrator/actions/x/reject",
	} {
		req := httptest.NewRequest(http.MethodPost, path, nil)
		rec := httptest.NewRecorder()
		NewRouter(Deps{}).ServeHTTP(rec, req)
		if rec.Code != http.StatusNotFound {
			t.Fatalf("%s status=%d", path, rec.Code)
		}
	}
}
