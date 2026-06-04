package clidevices

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"testing"
	"time"
)

func TestSignCSR_roundTrip(t *testing.T) {
	ca, _, _, err := GenerateDevCA()
	if err != nil {
		t.Fatal(err)
	}

	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	csrDER, err := x509.CreateCertificateRequest(rand.Reader, &x509.CertificateRequest{
		Subject: pkix.Name{CommonName: "cli-test-device"},
	}, key)
	if err != nil {
		t.Fatal(err)
	}
	csrPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csrDER})

	certPEM, cert, err := ca.SignCSR(csrPEM, 24*time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	if len(certPEM) == 0 {
		t.Fatal("certificado PEM vazio")
	}
	fp := Fingerprint(cert)
	if len(fp) != 32 {
		t.Fatalf("fingerprint len=%d", len(fp))
	}
}
