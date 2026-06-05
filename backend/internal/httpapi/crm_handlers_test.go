package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCRMRoutes_NotRegisteredWithoutDeps(t *testing.T) {
	routes := []struct {
		method string
		path   string
	}{
		{http.MethodGet, "/api/crm/leads"},
		{http.MethodPost, "/api/crm/leads"},
		{http.MethodPut, "/api/crm/leads/x"},
		{http.MethodDelete, "/api/crm/leads/x"},
	}
	for _, rt := range routes {
		req := httptest.NewRequest(rt.method, rt.path, nil)
		rec := httptest.NewRecorder()
		NewRouter(Deps{}).ServeHTTP(rec, req)
		if rec.Code != http.StatusNotFound {
			t.Fatalf("%s %s: status=%d", rt.method, rt.path, rec.Code)
		}
	}
}
