package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// Sem dependência HR (nem Auth), as rotas /api/hr/* não são registadas —
// devolvem 404, tal como as restantes rotas protegidas.
func TestEmployeeRecordRoutesGuarded(t *testing.T) {
	routes := []struct {
		method string
		path   string
	}{
		{http.MethodPost, "/api/hr/employees"},
		{http.MethodGet, "/api/hr/employees"},
		{http.MethodGet, "/api/hr/employees/abc"},
		{http.MethodDelete, "/api/hr/employees/abc"},
		{http.MethodPut, "/api/hr/employees/abc/fields/salary"},
		{http.MethodDelete, "/api/hr/employees/abc/fields/salary"},
		{http.MethodPost, "/api/hr/employees/abc/fields/salary/shred"},
		{http.MethodPost, "/api/hr/employees/abc/shred"},
		{http.MethodGet, "/api/hr/certificates"},
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
