package httpapi

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/burnnotes"
)

// Sem dependência BurnNotes, as rotas /api/share/notes* não são registadas.
func TestBurnNoteRoutesGuarded(t *testing.T) {
	routes := []struct {
		method string
		path   string
	}{
		{http.MethodPost, "/api/share/notes"},
		{http.MethodPost, "/api/share/notes/abc"},
		{http.MethodPost, "/api/share/notes/abc/burn"},
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

// A leitura é pública, devolve o ciphertext uma vez e queima a nota (2.ª = 404).
func TestConsumeBurnNotePublic(t *testing.T) {
	store := burnnotes.NewStore()
	id, _, _, err := store.Create([]byte("nota-cifrada"), time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	router := NewRouter(Deps{BurnNotes: store})

	do := func() *httptest.ResponseRecorder {
		req := httptest.NewRequest(http.MethodPost, "/api/share/notes/"+id, nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		return rec
	}

	rec := do()
	if rec.Code != http.StatusOK {
		t.Fatalf("primeira leitura: status=%d, esperado 200", rec.Code)
	}
	var body struct {
		Ciphertext string `json:"ciphertext"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if got, _ := base64.StdEncoding.DecodeString(body.Ciphertext); string(got) != "nota-cifrada" {
		t.Fatalf("ciphertext devolvido difere: %q", got)
	}
	if rec2 := do(); rec2.Code != http.StatusNotFound {
		t.Fatalf("segunda leitura: status=%d, esperado 404 (burn-after-read)", rec2.Code)
	}
}

// O burn manual exige o token: token errado dá 403, token certo destroi (204) e
// a leitura seguinte dá 404.
func TestBurnNoteManual(t *testing.T) {
	store := burnnotes.NewStore()
	id, token, _, err := store.Create([]byte("segredo"), time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	router := NewRouter(Deps{BurnNotes: store})

	burn := func(tok string) int {
		req := httptest.NewRequest(http.MethodPost, "/api/share/notes/"+id+"/burn",
			strings.NewReader(`{"burn_token":"`+tok+`"}`))
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		return rec.Code
	}

	if code := burn("token-errado"); code != http.StatusForbidden {
		t.Fatalf("burn token errado: status=%d, esperado 403", code)
	}
	if code := burn(token); code != http.StatusNoContent {
		t.Fatalf("burn token certo: status=%d, esperado 204", code)
	}
	// Já não há nada para ler.
	req := httptest.NewRequest(http.MethodPost, "/api/share/notes/"+id, nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("leitura apos burn: status=%d, esperado 404", rec.Code)
	}
}
