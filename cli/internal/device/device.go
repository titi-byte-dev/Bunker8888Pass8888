// Package device gera CSR e guarda certificados locais.
package device

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"os"

	clicfg "github.com/titi-byte-dev/Bunker8888Pass8888/cli/internal/config"
)

// GenerateCSR cria par de chaves ECDSA P-256 e CSR PEM.
func GenerateCSR(commonName string) (csrPEM string, keyPEM []byte, err error) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return "", nil, err
	}
	csrDER, err := x509.CreateCertificateRequest(rand.Reader, &x509.CertificateRequest{
		Subject: pkix.Name{CommonName: commonName},
	}, key)
	if err != nil {
		return "", nil, err
	}
	keyDER, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		return "", nil, err
	}
	return string(pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csrDER})),
		pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: keyDER}),
		nil
}

// SaveCredentials grava certificados com permissões restritas (0600).
func SaveCredentials(store *clicfg.Store, certPEM, caPEM, keyPEM []byte, email string) error {
	if err := store.EnsureDir(); err != nil {
		return err
	}
	if err := writeFile0600(store.ClientKey, keyPEM); err != nil {
		return err
	}
	if err := writeFile0600(store.ClientCert, certPEM); err != nil {
		return err
	}
	if err := writeFile0600(store.CACert, caPEM); err != nil {
		return err
	}
	return writeFile0600(store.Email, []byte(email))
}

func writeFile0600(path string, data []byte) error {
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("escrever %s: %w", path, err)
	}
	return nil
}

// ReadEmail devolve o email guardado no registo do dispositivo.
func ReadEmail(store *clicfg.Store) (string, error) {
	b, err := os.ReadFile(store.Email)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
