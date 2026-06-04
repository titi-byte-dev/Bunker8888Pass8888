package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

// Erros do pacote. Declará-los como variáveis permite a quem chama usar
// errors.Is(err, ErrInvalidKeySize) em vez de comparar strings.
var (
	ErrInvalidKeySize     = errors.New("crypto: a chave tem de ter 32 bytes (AES-256)")
	ErrCiphertextTooShort = errors.New("crypto: ciphertext demasiado curto para conter o nonce")
)

// Encrypt cifra `plaintext` com AES-256-GCM usando `key` (32 bytes).
//
// `aad` (Additional Authenticated Data) é opcional: dados que NÃO são cifrados
// mas cuja integridade é protegida (ex: um ID de item). Passar nil se não houver.
//
// O resultado tem o formato: nonce || ciphertext || tag, tudo junto. Guardamos
// o nonce com o ciphertext porque é preciso para decifrar e não é secreto.
//
// ⚠️ Segurança: AES-GCM dá confidencialidade E autenticidade. O nonce é gerado
// aleatoriamente a cada chamada e NUNCA deve ser reutilizado com a mesma chave.
func Encrypt(key, plaintext, aad []byte) ([]byte, error) {
	gcm, err := newGCM(key)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// Seal(dst, nonce, plaintext, aad): anexa o resultado a `dst`. Ao passar o
	// próprio `nonce` como dst, ficamos com nonce||ciphertext||tag num só slice.
	return gcm.Seal(nonce, nonce, plaintext, aad), nil
}

// Decrypt reverte Encrypt. Devolve erro se a autenticação falhar (dados
// adulterados, chave errada, ou `aad` diferente do usado a cifrar).
func Decrypt(key, blob, aad []byte) ([]byte, error) {
	gcm, err := newGCM(key)
	if err != nil {
		return nil, err
	}

	ns := gcm.NonceSize()
	if len(blob) < ns {
		return nil, ErrCiphertextTooShort
	}

	// Separamos o nonce (prefixo) do ciphertext+tag (resto).
	nonce, ciphertext := blob[:ns], blob[ns:]

	// Open verifica a tag de autenticação ANTES de devolver o plaintext. Se a
	// verificação falhar, devolve erro e NUNCA dados parciais.
	return gcm.Open(nil, nonce, ciphertext, aad)
}

// newGCM valida a chave e constrói o modo AES-GCM.
func newGCM(key []byte) (cipher.AEAD, error) {
	if len(key) != KeySize {
		return nil, ErrInvalidKeySize
	}
	block, err := aes.NewCipher(key) // AES-256 porque a chave tem 32 bytes
	if err != nil {
		return nil, err
	}
	return cipher.NewGCM(block)
}

// encryptWithNonce permite cifrar com um nonce explícito. É usado SÓ nos testes
// (para validar contra test vectors determinísticos). Não é exportado de propósito.
//
// ⚠️ Nunca usar isto em produção com um nonce fixo/reutilizado.
func encryptWithNonce(key, nonce, plaintext, aad []byte) ([]byte, error) {
	gcm, err := newGCM(key)
	if err != nil {
		return nil, err
	}
	return gcm.Seal(nil, nonce, plaintext, aad), nil
}
