package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// Sem dependência SharedVaults (nem Auth), as rotas /api/share/vaults* não são
// registadas — devolvem 404, tal como as restantes rotas protegidas.
func TestSharedVaultRoutesGuarded(t *testing.T) {
	routes := []struct {
		method string
		path   string
	}{
		{http.MethodPost, "/api/share/vaults"},
		{http.MethodGet, "/api/share/vaults"},
		{http.MethodGet, "/api/share/vaults/abc"},
		{http.MethodDelete, "/api/share/vaults/abc"},
		{http.MethodGet, "/api/share/vaults/abc/members"},
		{http.MethodPost, "/api/share/vaults/abc/members"},
		{http.MethodDelete, "/api/share/vaults/abc/members/u1"},
		{http.MethodGet, "/api/share/vaults/abc/items"},
		{http.MethodPost, "/api/share/vaults/abc/items"},
		{http.MethodDelete, "/api/share/vaults/abc/items/i1"},
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
