package googledrive_test

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/googledrive"
)

func testPool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	dsn := os.Getenv("AEGIS_TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("AEGIS_TEST_DATABASE_URL não definido")
	}
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		t.Fatal(err)
	}
	return pool
}

func TestRepo_CreateListDelete(t *testing.T) {
	pool := testPool(t)
	defer pool.Close()
	repo := googledrive.NewRepo(pool)
	ctx := context.Background()
	owner := "00000000-0000-0000-0000-000000000099"

	f, err := repo.Create(ctx, owner, []byte("opaque-blob"))
	if err != nil {
		t.Fatal(err)
	}
	list, err := repo.List(ctx, owner)
	if err != nil || len(list) == 0 {
		t.Fatalf("list=%v err=%v", list, err)
	}
	got, err := repo.Get(ctx, owner, f.ID)
	if err != nil || string(got.Blob) != "opaque-blob" {
		t.Fatalf("get=%+v err=%v", got, err)
	}
	if err := repo.Delete(ctx, owner, f.ID); err != nil {
		t.Fatal(err)
	}
}
