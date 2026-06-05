package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFinanceAgentRoute_NotRegisteredWithoutBus(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/agent/finance/report-stale", nil)
	rec := httptest.NewRecorder()
	NewRouter(Deps{}).ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status=%d, esperado 404", rec.Code)
	}
}
