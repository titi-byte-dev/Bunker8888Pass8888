package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestInventoryRoutesGuarded(t *testing.T) {
	routes := []struct {
		method string
		path   string
	}{
		{http.MethodGet, "/api/ops/inventory"},
		{http.MethodPost, "/api/ops/inventory"},
		{http.MethodPut, "/api/ops/inventory/x"},
		{http.MethodPost, "/api/ops/inventory/x/adjust"},
		{http.MethodDelete, "/api/ops/inventory/x"},
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
