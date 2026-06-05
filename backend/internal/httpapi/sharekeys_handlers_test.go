package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// Sem dependência ShareKeys (nem Auth), as rotas /api/share/* não são
// registadas — devolvem 404, tal como as restantes rotas protegidas.
func TestShareKeypairRoutesGuarded(t *testing.T) {
	routes := []struct {
		method string
		path   string
	}{
		{http.MethodGet, "/api/share/keypair"},
		{http.MethodPut, "/api/share/keypair"},
		{http.MethodGet, "/api/share/keypair/status"},
		{http.MethodGet, "/api/share/public-key?email=a@b.local"},
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
