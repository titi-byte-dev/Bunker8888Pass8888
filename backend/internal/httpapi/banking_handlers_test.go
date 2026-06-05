package httpapi_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/httpapi"
)

func TestBankingRoutes_NotRegisteredWithoutDeps(t *testing.T) {
	mux := httpapi.NewRouter(httpapi.Deps{})
	for _, path := range []string{
		"/api/fin/banking/status",
		"/api/fin/banking/connect",
		"/api/fin/banking/sync",
	} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, req)
		if rec.Code != http.StatusNotFound {
			t.Fatalf("%s: esperava 404, obteve %d", path, rec.Code)
		}
	}
}
