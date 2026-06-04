package auth

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/db"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/sessions"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/users"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/pkg/crypto"
)

func TestHashToken_Deterministic(t *testing.T) {
	a := hashToken("abc123")
	b := hashToken("abc123")
	if string(a) != string(b) {
		t.Fatal("hashToken não é determinístico")
	}
}

func TestHashToken_DifferentInputs(t *testing.T) {
	a := hashToken("token-a")
	b := hashToken("token-b")
	if string(a) == string(b) {
		t.Fatal("tokens diferentes produziram o mesmo hash")
	}
}

func TestNewToken_Unique(t *testing.T) {
	a, err := newToken()
	if err != nil {
		t.Fatal(err)
	}
	b, err := newToken()
	if err != nil {
		t.Fatal(err)
	}
	if a == b {
		t.Fatal("dois tokens aleatórios iguais")
	}
	if len(a) != 64 {
		t.Fatalf("token len=%d, esperado 64", len(a))
	}
}

func TestRegisterLoginFlow_Integration(t *testing.T) {
	url := os.Getenv("AEGIS_TEST_DATABASE_URL")
	if url == "" {
		t.Skip("AEGIS_TEST_DATABASE_URL não definido; a saltar integração")
	}

	ctx := context.Background()
	pool, err := db.Connect(ctx, url)
	if err != nil {
		t.Fatalf("ligar à BD: %v", err)
	}
	defer pool.Close()

	if err := db.Migrate(ctx, pool); err != nil {
		t.Fatalf("migrar: %v", err)
	}
	cleanTables(t, ctx, pool)

	userRepo := users.NewRepo(pool)
	sessionRepo := sessions.NewRepo(pool)
	svc := NewService(userRepo, sessionRepo, 3600)

	email := "test@example.com"
	authHash := []byte("auth-hash-simulado-do-cliente")
	kdf := ClientKDF{Salt: []byte("0123456789abcdef"), Time: 1, Memory: 8192, Threads: 1}

	if err := svc.Register(ctx, email, authHash, kdf); err != nil {
		t.Fatalf("Register: %v", err)
	}
	if err := svc.Register(ctx, email, authHash, kdf); err == nil {
		t.Fatal("registo duplicado devia falhar")
	}

	token, err := svc.Login(ctx, email, authHash)
	if err != nil {
		t.Fatalf("Login: %v", err)
	}
	userID, err := svc.Authenticate(ctx, token)
	if err != nil {
		t.Fatalf("Authenticate: %v", err)
	}
	if userID == "" {
		t.Fatal("userID vazio")
	}

	if _, err := svc.Login(ctx, email, []byte("hash-errado")); err == nil {
		t.Fatal("login com hash errado devia falhar")
	}

	if err := svc.Logout(ctx, token); err != nil {
		t.Fatalf("Logout: %v", err)
	}
	if _, err := svc.Authenticate(ctx, token); err == nil {
		t.Fatal("token após logout devia ser inválido")
	}

	u, _ := userRepo.ByEmail(ctx, email)
	if crypto.ConstantTimeEqual(u.Verifier, authHash) {
		t.Fatal("verifier igual ao auth hash cru")
	}
}

func cleanTables(t *testing.T, ctx context.Context, pool *pgxpool.Pool) {
	t.Helper()
	for _, q := range []string{
		`DELETE FROM vault_items`,
		`DELETE FROM sessions`,
		`DELETE FROM users`,
	} {
		if _, err := pool.Exec(ctx, q); err != nil {
			t.Fatalf("limpar BD: %v", err)
		}
	}
}
