package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// Sem dependência Fin (nem Auth), as rotas /api/fin/* não são registadas.
func TestSubscriptionRoutesGuarded(t *testing.T) {
	routes := []struct {
		method string
		path   string
	}{
		{http.MethodGet, "/api/fin/subscriptions"},
		{http.MethodPost, "/api/fin/subscriptions"},
		{http.MethodPut, "/api/fin/subscriptions/s1"},
		{http.MethodDelete, "/api/fin/subscriptions/s1"},
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
