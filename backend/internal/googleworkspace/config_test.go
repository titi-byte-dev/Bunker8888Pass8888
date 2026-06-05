package googleworkspace_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/googleworkspace"
)

func TestConfig_ProviderKind_Mock(t *testing.T) {
	cfg := googleworkspace.Config{Enabled: false}
	if cfg.ProviderKind() != "mock" {
		t.Fatalf("kind=%s", cfg.ProviderKind())
	}
}

func TestConfig_ProviderKind_ServiceAccount(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sa.json")
	if err := os.WriteFile(path, []byte(`{"type":"service_account"}`), 0o600); err != nil {
		t.Fatal(err)
	}
	cfg := googleworkspace.Config{
		Enabled:       true,
		SAJSONPath:    path,
		DelegatedUser: "admin@example.com",
	}
	if cfg.ProviderKind() != "service_account" {
		t.Fatalf("kind=%s", cfg.ProviderKind())
	}
}

func TestService_Status_Mock(t *testing.T) {
	cfg := googleworkspace.Config{}
	svc := googleworkspace.NewService(cfg, googleworkspace.MockProvider{})
	st := svc.Status(t.Context())
	if !st.Ready || st.Provider != "mock" {
		t.Fatalf("status=%+v", st)
	}
}
