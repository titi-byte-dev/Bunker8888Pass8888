package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecruitmentRoute_NotRegisteredWithoutDeps(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/agent/recruitment/run", nil)
	rec := httptest.NewRecorder()
	NewRouter(Deps{}).ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status=%d, esperado 404", rec.Code)
	}
}
