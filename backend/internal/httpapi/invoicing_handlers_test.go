package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// Sem dependência Invoicing (nem Auth), as rotas /api/fin/invoices/* não são registadas.
func TestInvoiceRoutesGuarded(t *testing.T) {
	routes := []struct {
		method string
		path   string
	}{
		{http.MethodGet, "/api/fin/invoices"},
		{http.MethodPost, "/api/fin/invoices"},
		{http.MethodPut, "/api/fin/invoices/i1/status"},
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
