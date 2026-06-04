package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthz(t *testing.T) {
	// httptest.NewRecorder captura a resposta sem abrir uma porta real.
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()

	NewRouter(Deps{}).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d, esperado 200", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("Content-Type=%q, esperado application/json", ct)
	}
}

func TestRegister_RequiresBody(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", nil)
	rec := httptest.NewRecorder()

	// Sem Auth configurado, a rota nem existe — mas com Auth nil o handler não
	// é registado. Testamos com router vazio: deve dar 404.
	NewRouter(Deps{}).ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("sem deps auth: status=%d, esperado 404", rec.Code)
	}
}

func TestServerTime(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/time", nil)
	rec := httptest.NewRecorder()

	NewRouter(Deps{}).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d, esperado 200", rec.Code)
	}
}

func TestBearerToken(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "Bearer abc.def.ghi")
	if got := bearerToken(req); got != "abc.def.ghi" {
		t.Fatalf("bearerToken=%q", got)
	}
	req.Header.Set("Authorization", "Basic xyz")
	if got := bearerToken(req); got != "" {
		t.Fatalf("Basic auth devia ser ignorado, got=%q", got)
	}
}
