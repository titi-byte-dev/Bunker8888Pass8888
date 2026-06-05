package httpapi

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/secretlinks"
)

// Sem dependência SecretLinks, as rotas /api/share/links* não são registadas.
func TestSecretLinkRoutesGuarded(t *testing.T) {
	routes := []struct {
		method string
		path   string
	}{
		{http.MethodPost, "/api/share/links"},
		{http.MethodPost, "/api/share/links/abc"},
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

// O consumo é público e devolve o ciphertext uma vez; à segunda dá 404.
func TestConsumeSecretLinkPublic(t *testing.T) {
	store := secretlinks.NewStore()
	id, _, err := store.Create([]byte("blob-cifrado"), time.Hour, 1)
	if err != nil {
		t.Fatal(err)
	}
	router := NewRouter(Deps{SecretLinks: store})

	do := func() *httptest.ResponseRecorder {
		req := httptest.NewRequest(http.MethodPost, "/api/share/links/"+id, nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		return rec
	}

	rec := do()
	if rec.Code != http.StatusOK {
		t.Fatalf("primeiro consumo: status=%d, esperado 200", rec.Code)
	}
	var body struct {
		Ciphertext string `json:"ciphertext"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if got, _ := base64.StdEncoding.DecodeString(body.Ciphertext); string(got) != "blob-cifrado" {
		t.Fatalf("ciphertext devolvido difere: %q", got)
	}

	if rec2 := do(); rec2.Code != http.StatusNotFound {
		t.Fatalf("segundo consumo: status=%d, esperado 404 (uso unico)", rec2.Code)
	}
}
