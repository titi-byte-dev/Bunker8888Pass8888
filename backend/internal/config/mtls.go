package config

import (
	"fmt"
	"os"

	"github.com/titi-byte-dev/Bunker8888Pass8888/backend/internal/clidevices"
)

// MTLSMaterial agrupa certificados para o servidor mTLS da CLI.
type MTLSMaterial struct {
	CA           *clidevices.CA
	CACertPEM    []byte
	ServerCertPEM []byte
	ServerKeyPEM  []byte
}

// LoadMTLSMaterial carrega ou gera (dev) os certificados mTLS.
func LoadMTLSMaterial(cfg Config) (*MTLSMaterial, error) {
	if cfg.MTLSAutoDev {
		ca, caCert, _, err := clidevices.GenerateDevCA()
		if err != nil {
			return nil, err
		}
		srvCert, srvKey, err := ca.ServerTLSCert("aegis-mtls")
		if err != nil {
			return nil, err
		}
		return &MTLSMaterial{
			CA:            ca,
			CACertPEM:     caCert,
			ServerCertPEM: srvCert,
			ServerKeyPEM:  srvKey,
		}, nil
	}
	if cfg.MTLSCACert == "" || cfg.MTLSCAKey == "" || cfg.MTLSServerCert == "" || cfg.MTLSServerKey == "" {
		return nil, nil
	}
	caCert, err := os.ReadFile(cfg.MTLSCACert)
	if err != nil {
		return nil, fmt.Errorf("ler CA cert: %w", err)
	}
	caKey, err := os.ReadFile(cfg.MTLSCAKey)
	if err != nil {
		return nil, fmt.Errorf("ler CA key: %w", err)
	}
	srvCert, err := os.ReadFile(cfg.MTLSServerCert)
	if err != nil {
		return nil, fmt.Errorf("ler server cert: %w", err)
	}
	srvKey, err := os.ReadFile(cfg.MTLSServerKey)
	if err != nil {
		return nil, fmt.Errorf("ler server key: %w", err)
	}
	ca, err := clidevices.LoadCA(caCert, caKey)
	if err != nil {
		return nil, err
	}
	return &MTLSMaterial{
		CA:            ca,
		CACertPEM:     caCert,
		ServerCertPEM: srvCert,
		ServerKeyPEM:  srvKey,
	}, nil
}
