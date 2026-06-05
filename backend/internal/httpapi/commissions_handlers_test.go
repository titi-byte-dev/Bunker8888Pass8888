package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// Sem dependência Commissions (nem Auth), as rotas /api/fin/commissions/* não são registadas.
func TestCommissionRoutesGuarded(t *testing.T) {
	routes := []struct {
		method string
		path   string
	}{
		{http.MethodGet, "/api/fin/commissions"},
		{http.MethodPost, "/api/fin/commissions"},
		{http.MethodPut, "/api/fin/commissions/c1/status"},
	}
	for _, rt := range routes {
		req := httptest.NewRequest(rt.method, rt.path, nil)
		rec := httptest.NewRecorder()
		NewRouter(Deps{}).ServeHTTP(rec, req)
		if rec.Code != http.StatusNotFound {
			t.Fatalf("%s %s: status=%d, esperado 404 sem deps", rt.method, rt.path, rec.Code)
		}
	}
}
