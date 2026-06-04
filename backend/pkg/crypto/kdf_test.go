package crypto

import (
	"bytes"
	"encoding/hex"
	"testing"
)

// testParams usa custos baixos para os testes correrem depressa. Em produção
// usam-se os DefaultKDFParams (muito mais pesados).
func testParams() KDFParams {
	return KDFParams{TimeCost: 1, MemoryKiB: 8 * 1024, Threads: 1, KeyLen: KeySize}
}

func TestDeriveMasterKey_Deterministic(t *testing.T) {
	// Mesma password + mesmo salt + mesmos params => SEMPRE a mesma chave.
	// É isto que permite ao utilizador voltar a abrir o cofre noutro dispositivo.
	pw := []byte("password-super-secreta")
	salt := []byte("0123456789abcdef")

	k1 := DeriveMasterKey(pw, salt, testParams())
	k2 := DeriveMasterKey(pw, salt, testParams())

	if !bytes.Equal(k1, k2) {
		t.Fatalf("derivação não determinística: %x != %x", k1, k2)
	}
	if len(k1) != KeySize {
		t.Fatalf("tamanho da chave = %d, esperado %d", len(k1), KeySize)
	}
}

func TestDeriveMasterKey_SaltMatters(t *testing.T) {
	// Salts diferentes têm de produzir chaves diferentes (senão o salt seria inútil).
	pw := []byte("a-mesma-password")
	k1 := DeriveMasterKey(pw, []byte("salt-aaaaaaaaaaa"), testParams())
	k2 := DeriveMasterKey(pw, []byte("salt-bbbbbbbbbbb"), testParams())

	if bytes.Equal(k1, k2) {
		t.Fatal("salts diferentes produziram a mesma chave")
	}
}

func TestGenerateSalt(t *testing.T) {
	s1, err := GenerateSalt()
	if err != nil {
		t.Fatalf("GenerateSalt: %v", err)
	}
	if len(s1) != SaltLen {
		t.Fatalf("tamanho do salt = %d, esperado %d", len(s1), SaltLen)
	}
	// Dois salts seguidos não devem ser iguais (probabilidade desprezável).
	s2, _ := GenerateSalt()
	if bytes.Equal(s1, s2) {
		t.Fatal("dois salts aleatórios saíram iguais")
	}
}

func TestDeriveAuthHash_DiffersFromMasterKey(t *testing.T) {
	// ⚠️ O auth hash (enviado ao servidor) NÃO pode ser igual à master key
	// (que fica no cliente). São derivações independentes.
	pw := []byte("password")
	salt := []byte("0123456789abcdef")
	mk := DeriveMasterKey(pw, salt, testParams())
	auth := DeriveAuthHash(mk, pw, testParams())

	if bytes.Equal(mk, auth) {
		t.Fatal("auth hash igual à master key — quebraria o Zero-Knowledge")
	}
}

func TestConstantTimeEqual(t *testing.T) {
	a := []byte("abcdef")
	if !ConstantTimeEqual(a, []byte("abcdef")) {
		t.Fatal("valores iguais reportados como diferentes")
	}
	if ConstantTimeEqual(a, []byte("abcdeg")) {
		t.Fatal("valores diferentes reportados como iguais")
	}
}

// TestDeriveMasterKey_RegressionVector fixa um valor conhecido para detetar
// alterações acidentais no algoritmo/parâmetros. NÃO é um KAT do RFC 9106
// (a API do x/crypto não expõe os parâmetros 'secret'/'associated data' do RFC),
// mas garante estabilidade da nossa derivação ao longo do tempo.
func TestDeriveMasterKey_RegressionVector(t *testing.T) {
	pw := []byte("correct horse battery staple")
	salt := []byte("aegis-fixed-salt") // 16 bytes
	got := hex.EncodeToString(DeriveMasterKey(pw, salt, testParams()))

	// Valor fixado após a primeira execução verificada (params de teste).
	const want = "5a4f94c2f369326971d3753d5d7960dd6392f79060136b2b46a85d1bae3499e1"
	if got != want {
		t.Fatalf("vetor de regressão falhou:\n got=%s\nwant=%s", got, want)
	}
}
