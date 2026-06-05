package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMailpitWebhook_NotRegisteredWithoutDeps(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/mail/webhook/mailpit?secret=x", nil)
	rec := httptest.NewRecorder()
	NewRouter(Deps{}).ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status=%d", rec.Code)
	}
}
