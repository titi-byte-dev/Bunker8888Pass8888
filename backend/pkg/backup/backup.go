// Package backup cifra/decifra dumps da base de dados (INFRA-004).
//
// Didático: os backups em disco são um alvo frequente de exfiltração. Ciframos
// o dump com AES-256-GCM antes de o escrever — mesmo que alguém copie o ficheiro
// .enc, não consegue ler sem a chave AEGIS_BACKUP_KEY.
package backup

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"

	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/pkg/crypto"
)

// Magic identifica ficheiros de backup AegisPass (8 bytes ASCII).
var Magic = []byte("AEGISBK1")

var (
	ErrInvalidKey  = errors.New("backup: chave inválida (esperado 32 bytes em base64)")
	ErrInvalidFile = errors.New("backup: ficheiro inválido ou corrompido")
	ErrKeyNotSet   = errors.New("backup: AEGIS_BACKUP_KEY não definida")
)

// KeyFromEnv lê AEGIS_BACKUP_KEY (base64, 32 bytes decodificados).
func KeyFromEnv() ([]byte, error) {
	raw := os.Getenv("AEGIS_BACKUP_KEY")
	if raw == "" {
		return nil, ErrKeyNotSet
	}
	key, err := base64.StdEncoding.DecodeString(raw)
	if err != nil || len(key) != crypto.KeySize {
		return nil, ErrInvalidKey
	}
	return key, nil
}

// EncryptDump cifra plaintext e prefixa o cabeçalho mágico.
func EncryptDump(key, plaintext []byte) ([]byte, error) {
	blob, err := crypto.Encrypt(key, plaintext, Magic)
	if err != nil {
		return nil, err
	}
	out := make([]byte, len(Magic)+len(blob))
	copy(out, Magic)
	copy(out[len(Magic):], blob)
	return out, nil
}

// DecryptDump reverte EncryptDump.
func DecryptDump(key, file []byte) ([]byte, error) {
	if len(file) < len(Magic) || string(file[:len(Magic)]) != string(Magic) {
		return nil, ErrInvalidFile
	}
	return crypto.Decrypt(key, file[len(Magic):], Magic)
}

// GenerateKeyBase64 produz uma chave aleatória codificada em base64 (para .env).
func GenerateKeyBase64() (string, error) {
	salt, err := crypto.GenerateSalt()
	if err != nil {
		return "", err
	}
	extra, err := crypto.GenerateSalt()
	if err != nil {
		return "", err
	}
	key := append(salt, extra...)
	if len(key) != crypto.KeySize {
		return "", fmt.Errorf("tamanho inesperado da chave")
	}
	return base64.StdEncoding.EncodeToString(key), nil
}
