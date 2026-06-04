package backup

import (
	"bytes"
	"os"
	"testing"
)

func TestEncryptDecryptDump_RoundTrip(t *testing.T) {
	key := bytes.Repeat([]byte{0xAB}, 32)
	plain := []byte("-- PostgreSQL dump simulado\nCREATE TABLE users;\n")
	enc, err := EncryptDump(key, plain)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.HasPrefix(enc, Magic) {
		t.Fatal("falta cabeçalho mágico")
	}
	got, err := DecryptDump(key, enc)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, plain) {
		t.Fatal("round-trip falhou")
	}
}

func TestDecryptDump_WrongKey(t *testing.T) {
	key := bytes.Repeat([]byte{0x01}, 32)
	enc, _ := EncryptDump(key, []byte("x"))
	wrong := bytes.Repeat([]byte{0x02}, 32)
	if _, err := DecryptDump(wrong, enc); err == nil {
		t.Fatal("chave errada devia falhar")
	}
}

func TestGenerateKeyBase64(t *testing.T) {
	b64, err := GenerateKeyBase64()
	if err != nil {
		t.Fatal(err)
	}
	os.Setenv("AEGIS_BACKUP_KEY", b64)
	t.Cleanup(func() { os.Unsetenv("AEGIS_BACKUP_KEY") })
	key, err := KeyFromEnv()
	if err != nil || len(key) != 32 {
		t.Fatalf("KeyFromEnv: err=%v len=%d", err, len(key))
	}
}
