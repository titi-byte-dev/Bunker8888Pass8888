package crypto

import (
	"crypto/rand"
	"crypto/subtle"
	"io"

	"golang.org/x/crypto/argon2"
)

// KeySize é o tamanho, em bytes, das chaves simétricas do AegisPass.
// 32 bytes = 256 bits → alimenta tanto a saída da KDF como o AES-256-GCM.
const KeySize = 32

// KDFParams controla o custo da derivação de chave com Argon2id.
//
// Didático: uma KDF (Key Derivation Function) transforma uma password (fraca,
// escolhida por humanos) numa chave criptográfica (forte). Argon2id é
// deliberadamente LENTA e usa muita memória, para tornar ataques de força bruta
// caríssimos. Os parâmetros são um compromisso entre segurança e tempo de espera.
type KDFParams struct {
	TimeCost  uint32 // nº de passagens (iterações)
	MemoryKiB uint32 // memória usada, em KiB
	Threads   uint8  // paralelismo (nº de lanes)
	KeyLen    uint32 // tamanho da chave de saída, em bytes
}

// DefaultKDFParams devolve parâmetros adequados a produção (~64 MiB de memória).
//
// ⚠️ Segurança: estes valores devem ser revistos periodicamente à medida que o
// hardware evolui. O objetivo é ~0.5–1s de derivação no dispositivo do utilizador.
func DefaultKDFParams() KDFParams {
	return KDFParams{
		TimeCost:  3,
		MemoryKiB: 64 * 1024, // 64 MiB
		Threads:   4,
		KeyLen:    KeySize, // 32 bytes → chave AES-256
	}
}

// SaltLen é o tamanho recomendado do salt, em bytes.
const SaltLen = 16

// GenerateSalt produz um salt aleatório criptograficamente seguro.
//
// ⚠️ Segurança: usamos crypto/rand (e NUNCA math/rand). O salt não precisa de ser
// secreto, mas tem de ser único por utilizador para evitar ataques com tabelas
// pré-computadas (rainbow tables).
func GenerateSalt() ([]byte, error) {
	salt := make([]byte, SaltLen)
	// io.ReadFull garante que enchemos TODOS os bytes do slice (ou devolve erro).
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, err
	}
	return salt, nil
}

// DeriveMasterKey deriva a Master Key a partir da password e do salt.
//
// No modelo Zero-Knowledge esta chave NUNCA é enviada ao servidor: serve para
// cifrar/decifrar os itens do cofre localmente.
func DeriveMasterKey(password, salt []byte, p KDFParams) []byte {
	return argon2.IDKey(password, salt, p.TimeCost, p.MemoryKiB, p.Threads, p.KeyLen)
}

// DeriveAuthHash deriva, a partir da Master Key, um valor de autenticação que
// PODE ser enviado ao servidor para provar a posse da password.
//
// ⚠️ Segurança: é uma SEGUNDA derivação, independente. O servidor guarda este
// auth hash — que não decifra nada. Enviar a própria Master Key destruiria o
// modelo Zero-Knowledge. Usamos a password como salt desta segunda passagem
// para que o resultado seja distinto da Master Key.
func DeriveAuthHash(masterKey, password []byte, p KDFParams) []byte {
	return argon2.IDKey(masterKey, password, 1, p.MemoryKiB, p.Threads, p.KeyLen)
}

// ConstantTimeEqual compara dois valores em tempo constante.
//
// ⚠️ Segurança: comparar hashes com "==" ou bytes.Equal pode revelar, pelo tempo
// de execução, quantos bytes coincidiram (timing attack). subtle.ConstantTimeCompare
// demora sempre o mesmo, independentemente de onde está a diferença.
func ConstantTimeEqual(a, b []byte) bool {
	return subtle.ConstantTimeCompare(a, b) == 1
}
