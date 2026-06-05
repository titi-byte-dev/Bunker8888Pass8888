package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequireAdminKey(t *testing.T) {
	called := false
	h := requireAdminKey("secret", func(w http.ResponseWriter, _ *http.Request) {
		called = true
		w.WriteHeader(http.StatusNoContent)
	})

	t.Run("sem chave", func(t *testing.T) {
		called = false
		req := httptest.NewRequest(http.MethodGet, "/api/admin/users", nil)
		rec := httptest.NewRecorder()
		h(rec, req)
		if rec.Code != http.StatusForbidden {
			t.Fatalf("status=%d", rec.Code)
		}
		if called {
			t.Fatal("handler não devia correr")
		}
	})

	t.Run("chave válida", func(t *testing.T) {
		called = false
		req := httptest.NewRequest(http.MethodGet, "/api/admin/users", nil)
		req.Header.Set("X-Admin-Key", "secret")
		rec := httptest.NewRecorder()
		h(rec, req)
		if rec.Code != http.StatusNoContent {
			t.Fatalf("status=%d", rec.Code)
		}
		if !called {
			t.Fatal("handler devia correr")
		}
	})

	t.Run("admin desactivado", func(t *testing.T) {
		hOff := requireAdminKey("", func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
		})
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("X-Admin-Key", "secret")
		rec := httptest.NewRecorder()
		hOff(rec, req)
		if rec.Code != http.StatusServiceUnavailable {
			t.Fatalf("status=%d", rec.Code)
		}
	})
}

func TestRegisterAdminRoutes_DisabledWithoutKey(t *testing.T) {
	mux := http.NewServeMux()
	registerAdminRoutes(mux, Deps{AdminKey: ""})
	req := httptest.NewRequest(http.MethodGet, "/api/admin/users", nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("sem admin key: status=%d, esperado 404", rec.Code)
	}
}
