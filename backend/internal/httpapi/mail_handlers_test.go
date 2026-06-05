package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// Sem dependência Mail (nem Auth), as rotas /api/mail/* não são registadas.
func TestMailAliasRoutesGuarded(t *testing.T) {
	routes := []struct {
		method string
		path   string
	}{
		{http.MethodGet, "/api/mail/aliases"},
		{http.MethodPost, "/api/mail/aliases"},
		{http.MethodPatch, "/api/mail/aliases/a1"},
		{http.MethodDelete, "/api/mail/aliases/a1"},
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
