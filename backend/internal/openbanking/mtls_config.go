package openbanking

import (
	"errors"
	"os"
)

// MTLSConfig reúne os caminhos PEM para falar com um TPP/banco real (FIN-003).
// Didático: Open Banking na UE exige mTLS — o cliente (nós) apresenta
// certificado à instituição financeira, além do TLS normal do servidor.
type MTLSConfig struct {
	ClientCertPath string
	ClientKeyPath  string
	CACertPath     string
	BaseURL        string
}

// LoadMTLSConfigFromEnv lê AEGIS_OB_* — vazio significa «só mock provider».
func LoadMTLSConfigFromEnv() (MTLSConfig, error) {
	cfg := MTLSConfig{
		ClientCertPath: os.Getenv("AEGIS_OB_CLIENT_CERT"),
		ClientKeyPath:  os.Getenv("AEGIS_OB_CLIENT_KEY"),
		CACertPath:     os.Getenv("AEGIS_OB_CA_CERT"),
		BaseURL:        os.Getenv("AEGIS_OB_BASE_URL"),
	}
	if cfg.BaseURL == "" {
		return cfg, nil
	}
	if cfg.ClientCertPath == "" || cfg.ClientKeyPath == "" {
		return MTLSConfig{}, errors.New("openbanking: AEGIS_OB_BASE_URL definido mas faltam certificados cliente")
	}
	return cfg, nil
}
