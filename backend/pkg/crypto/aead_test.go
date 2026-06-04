package crypto

import (
	"bytes"
	"encoding/hex"
	"errors"
	"testing"
)

func newTestKey(t *testing.T) []byte {
	t.Helper()
	k, err := GenerateSalt() // 16 bytes; precisamos de 32 → derivamos
	if err != nil {
		t.Fatal(err)
	}
	return DeriveMasterKey([]byte("k"), k, testParams())
}

func TestEncryptDecrypt_RoundTrip(t *testing.T) {
	key := newTestKey(t)
	casos := []struct {
		nome      string
		plaintext []byte
		aad       []byte
	}{
		{"texto simples", []byte("ola mundo"), nil},
		{"vazio", []byte(""), nil},
		{"binario", []byte{0x00, 0xFF, 0x10, 0x7F}, nil},
		{"com aad", []byte("dados"), []byte("item-id-123")},
	}
	for _, c := range casos {
		t.Run(c.nome, func(t *testing.T) {
			blob, err := Encrypt(key, c.plaintext, c.aad)
			if err != nil {
				t.Fatalf("Encrypt: %v", err)
			}
			got, err := Decrypt(key, blob, c.aad)
			if err != nil {
				t.Fatalf("Decrypt: %v", err)
			}
			if !bytes.Equal(got, c.plaintext) {
				t.Fatalf("round-trip falhou: got %q want %q", got, c.plaintext)
			}
		})
	}
}

func TestDecrypt_TamperDetected(t *testing.T) {
	key := newTestKey(t)
	blob, _ := Encrypt(key, []byte("mensagem secreta"), nil)

	// Alterar um byte do ciphertext deve fazer a autenticação (tag GCM) falhar.
	blob[len(blob)-1] ^= 0x01
	if _, err := Decrypt(key, blob, nil); err == nil {
		t.Fatal("adulteração não detetada — a autenticação GCM falhou em proteger")
	}
}

func TestDecrypt_WrongAAD(t *testing.T) {
	key := newTestKey(t)
	blob, _ := Encrypt(key, []byte("x"), []byte("aad-original"))
	if _, err := Decrypt(key, blob, []byte("aad-diferente")); err == nil {
		t.Fatal("AAD diferente devia falhar a autenticação")
	}
}

func TestDecrypt_WrongKey(t *testing.T) {
	blob, _ := Encrypt(newTestKey(t), []byte("x"), nil)
	if _, err := Decrypt(newTestKey(t), blob, nil); err == nil {
		t.Fatal("chave errada devia falhar")
	}
}

func TestDecrypt_TooShort(t *testing.T) {
	key := newTestKey(t)
	if _, err := Decrypt(key, []byte{0x00}, nil); !errors.Is(err, ErrCiphertextTooShort) {
		t.Fatalf("esperado ErrCiphertextTooShort, obtido %v", err)
	}
}

func TestInvalidKeySize(t *testing.T) {
	if _, err := Encrypt([]byte("curta"), []byte("x"), nil); !errors.Is(err, ErrInvalidKeySize) {
		t.Fatalf("esperado ErrInvalidKeySize, obtido %v", err)
	}
}

func TestEncrypt_NonceUniqueness(t *testing.T) {
	key := newTestKey(t)
	// Cifrar o MESMO texto duas vezes deve dar blobs diferentes, porque o nonce
	// é aleatório. Se fossem iguais, o nonce estaria a ser reutilizado (grave).
	a, _ := Encrypt(key, []byte("igual"), nil)
	b, _ := Encrypt(key, []byte("igual"), nil)
	if bytes.Equal(a, b) {
		t.Fatal("dois ciphertexts iguais — possível reutilização de nonce")
	}
}

// TestKAT_AES256GCM valida a nossa implementação contra um Known Answer Test
// publicado: "The Galois/Counter Mode of Operation (GCM)", McGrew & Viega,
// Test Case 16 (AES-256, com AAD). É a prova de correção independente exigida
// por VAULT-002 (test vectors).
func TestKAT_AES256GCM(t *testing.T) {
	key := mustHex(t, "feffe9928665731c6d6a8f9467308308feffe9928665731c6d6a8f9467308308")
	nonce := mustHex(t, "cafebabefacedbaddecaf888")
	plaintext := mustHex(t, "d9313225f88406e5a55909c5aff5269a86a7a9531534f7da2e4c303d8a318a721c3c0c95956809532fcf0e2449a6b525b16aedf5aa0de657ba637b39")
	aad := mustHex(t, "feedfacedeadbeeffeedfacedeadbeefabaddad2")
	wantCipher := "522dc1f099567d07f47f37a32a84427d643a8cdcbfe5c0c97598a2bd2555d1aa8cb08e48590dbb3da7b08b1056828838c5f61e6393ba7a0abcc9f662"
	wantTag := "76fc6ece0f4e1768cddf8853bb2d551b"

	out, err := encryptWithNonce(key, nonce, plaintext, aad)
	if err != nil {
		t.Fatalf("encryptWithNonce: %v", err)
	}
	if got := hex.EncodeToString(out); got != wantCipher+wantTag {
		t.Fatalf("KAT falhou:\n got=%s\nwant=%s", got, wantCipher+wantTag)
	}
}

func mustHex(t *testing.T, s string) []byte {
	t.Helper()
	b, err := hex.DecodeString(s)
	if err != nil {
		t.Fatalf("hex inválido: %v", err)
	}
	return b
}
