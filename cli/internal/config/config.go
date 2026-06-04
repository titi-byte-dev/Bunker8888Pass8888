// Package config gere ficheiros locais da CLI (~/.aegis ou AEGIS_CONFIG_DIR).
package config

import (
	"os"
	"path/filepath"
)

// Store paths para certificados e metadados da CLI.
type Store struct {
	Dir        string
	ClientCert string
	ClientKey  string
	CACert     string
	Email      string
}

// DefaultAPIBase é o endpoint HTTP da API (login, registo de dispositivo).
const DefaultAPIBase = "http://localhost:8080"

// DefaultMTLSBase é o endpoint mTLS para operações de cofre.
const DefaultMTLSBase = "https://localhost:8443"

// LoadStore resolve o directório de configuração.
func LoadStore() (*Store, error) {
	dir := os.Getenv("AEGIS_CONFIG_DIR")
	if dir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		dir = filepath.Join(home, ".aegis")
	}
	return &Store{
		Dir:        dir,
		ClientCert: filepath.Join(dir, "client.crt"),
		ClientKey:  filepath.Join(dir, "client.key"),
		CACert:     filepath.Join(dir, "ca.crt"),
		Email:      filepath.Join(dir, "email"),
	}, nil
}

// APIBase devolve a URL base HTTP (AEGIS_API_URL ou omissão).
func APIBase() string {
	if v := os.Getenv("AEGIS_API_URL"); v != "" {
		return v
	}
	return DefaultAPIBase
}

// MTLSBase devolve a URL base mTLS (AEGIS_MTLS_URL ou omissão).
func MTLSBase() string {
	if v := os.Getenv("AEGIS_MTLS_URL"); v != "" {
		return v
	}
	return DefaultMTLSBase
}

// EnsureDir cria o directório de config com permissões restritas (0700).
func (s *Store) EnsureDir() error {
	return os.MkdirAll(s.Dir, 0o700)
}

// HasDevice indica se certificados cliente estão presentes.
func (s *Store) HasDevice() bool {
	if _, err := os.Stat(s.ClientCert); err != nil {
		return false
	}
	if _, err := os.Stat(s.ClientKey); err != nil {
		return false
	}
	return true
}
