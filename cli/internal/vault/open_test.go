package vault

import (
	"encoding/base64"
	"encoding/json"
	"testing"

	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/pkg/crypto"
)

func TestOpenLogin_roundTrip(t *testing.T) {
	key := make([]byte, 32)
	for i := range key {
		key[i] = byte(i)
	}
	inner, err := json.Marshal(LoginItem{
		Kind: "login", Title: "DB", Username: "admin", Password: "test-value-xyz",
	})
	if err != nil {
		t.Fatal(err)
	}
	env, err := json.Marshal(map[string]any{
		"v": 1, "kind": "login", "data": json.RawMessage(inner),
	})
	if err != nil {
		t.Fatal(err)
	}
	blob, err := crypto.Encrypt(key, env, nil)
	if err != nil {
		t.Fatal(err)
	}
	b64 := base64.StdEncoding.EncodeToString(blob)
	item, err := OpenLogin(key, []byte(b64))
	if err != nil {
		t.Fatal(err)
	}
	if item.Password != "test-value-xyz" {
		t.Fatalf("password=%q", item.Password)
	}
	v, err := FieldValue(item, "password")
	if err != nil || v != "test-value-xyz" {
		t.Fatalf("field=%q err=%v", v, err)
	}
}
