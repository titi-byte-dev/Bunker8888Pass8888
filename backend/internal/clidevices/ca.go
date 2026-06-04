// Package clidevices gere dispositivos CLI autenticados por mTLS (VAULT-017).
//
// > 💡 **Conceito — mTLS:** no TLS normal só o servidor prova identidade. No mTLS
// ambos apresentam certificado — a CLI prova que é um dispositivo registado.
package clidevices

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"fmt"
	"math/big"
	"time"
)

var (
	ErrInvalidCSR     = errors.New("clidevices: CSR inválido")
	ErrUnsupportedKey = errors.New("clidevices: chave do CSR não suportada (use ECDSA P-256)")
)

// CA assina certificados cliente para dispositivos CLI.
type CA struct {
	cert *x509.Certificate
	key  *ecdsa.PrivateKey
}

// LoadCA carrega certificado e chave PEM da autoridade certificadora interna.
func LoadCA(certPEM, keyPEM []byte) (*CA, error) {
	cert, err := parseCertPEM(certPEM)
	if err != nil {
		return nil, err
	}
	key, err := parseECKeyPEM(keyPEM)
	if err != nil {
		return nil, err
	}
	return &CA{cert: cert, key: key}, nil
}

// GenerateDevCA cria uma CA efémera para desenvolvimento local.
//
// ⚠️ **Segurança:** só para dev — nunca usar CA gerada em runtime em produção.
func GenerateDevCA() (*CA, []byte, []byte, error) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, nil, err
	}
	tmpl := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{CommonName: "AegisPass Dev CA"},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(10 * 365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, tmpl, &key.PublicKey, key)
	if err != nil {
		return nil, nil, nil, err
	}
	cert, err := x509.ParseCertificate(der)
	if err != nil {
		return nil, nil, nil, err
	}
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})
	keyDER, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		return nil, nil, nil, err
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: keyDER})
	return &CA{cert: cert, key: key}, certPEM, keyPEM, nil
}

// SignCSR assina um pedido de certificado (CSR) e devolve o certificado cliente PEM.
func (ca *CA) SignCSR(csrPEM []byte, ttl time.Duration) ([]byte, *x509.Certificate, error) {
	block, _ := pem.Decode(csrPEM)
	if block == nil || block.Type != "CERTIFICATE REQUEST" {
		return nil, nil, ErrInvalidCSR
	}
	csr, err := x509.ParseCertificateRequest(block.Bytes)
	if err != nil {
		return nil, nil, ErrInvalidCSR
	}
	if err := csr.CheckSignature(); err != nil {
		return nil, nil, ErrInvalidCSR
	}

	pub, ok := csr.PublicKey.(*ecdsa.PublicKey)
	if !ok || pub.Curve != elliptic.P256() {
		return nil, nil, ErrUnsupportedKey
	}

	serial, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return nil, nil, err
	}
	now := time.Now()
	tmpl := &x509.Certificate{
		SerialNumber: serial,
		Subject:      csr.Subject,
		NotBefore:    now.Add(-time.Minute),
		NotAfter:     now.Add(ttl),
		KeyUsage:     x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, ca.cert, csr.PublicKey, ca.key)
	if err != nil {
		return nil, nil, err
	}
	clientCert, err := x509.ParseCertificate(der)
	if err != nil {
		return nil, nil, err
	}
	out := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})
	return out, clientCert, nil
}

// CertDER devolve o certificado da CA em DER (para incluir nas respostas de registo).
func (ca *CA) CertDER() []byte {
	return ca.cert.Raw
}

// Fingerprint devolve SHA-256 do certificado (identificador único do dispositivo).
func Fingerprint(cert *x509.Certificate) []byte {
	sum := sha256.Sum256(cert.Raw)
	return sum[:]
}

// ServerTLSCert gera um par servidor TLS assinado pela CA (dev).
func (ca *CA) ServerTLSCert(host string) ([]byte, []byte, error) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	serial, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return nil, nil, err
	}
	tmpl := &x509.Certificate{
		SerialNumber: serial,
		Subject:      pkix.Name{CommonName: host},
		NotBefore:    time.Now().Add(-time.Hour),
		NotAfter:     time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		DNSNames:     []string{host, "localhost"},
	}
	der, err := x509.CreateCertificate(rand.Reader, tmpl, ca.cert, &key.PublicKey, ca.key)
	if err != nil {
		return nil, nil, err
	}
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: der})
	keyDER, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		return nil, nil, err
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: keyDER})
	return certPEM, keyPEM, nil
}

func parseCertPEM(pemBytes []byte) (*x509.Certificate, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, fmt.Errorf("clidevices: PEM de certificado inválido")
	}
	return x509.ParseCertificate(block.Bytes)
}

func parseECKeyPEM(pemBytes []byte) (*ecdsa.PrivateKey, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, fmt.Errorf("clidevices: PEM de chave inválido")
	}
	return x509.ParseECPrivateKey(block.Bytes)
}
