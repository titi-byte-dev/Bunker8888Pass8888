// Package config carrega a configuração da aplicação a partir do ambiente.
//
// Didático: em Go é idiomático ter um pacote pequeno e focado para a config.
// Lemos de variáveis de ambiente (12-factor app) com valores por omissão
// sensatos, para que a app arranque em desenvolvimento sem configuração extra.
package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// Config agrupa todas as definições de runtime do servidor.
type Config struct {
	// HTTPAddr é o endereço onde a API escuta (ex: ":8080").
	HTTPAddr string
	// ShutdownTimeout é quanto tempo damos aos pedidos em curso para terminarem
	// quando o servidor recebe um sinal de paragem (graceful shutdown).
	ShutdownTimeout time.Duration
	// DatabaseURL é a string de ligação ao PostgreSQL (DSN).
	DatabaseURL string
	// AdminKey protege endpoints administrativos (remote wipe de terceiros).
	// ⚠️ Segurança: vazio desactiva esses endpoints — nunca usar valor por omissão.
	AdminKey string
	// MTLSAddr é o endereço do servidor mTLS da CLI (ex: ":8443"). Vazio desactiva.
	MTLSAddr string
	// MTLSAutoDev gera CA/servidor efémeros ao arrancar (só desenvolvimento).
	MTLSAutoDev bool
	// Caminhos PEM para mTLS em produção.
	MTLSCACert     string
	MTLSCAKey      string
	MTLSServerCert string
	MTLSServerKey  string
	// WebAuthn / passkeys (VAULT-014)
	WebAuthnRPDisplayName string
	WebAuthnRPID          string
	WebAuthnRPOrigins     []string
	// MAIL-002 (dev): Mailpit para SMTP de teste e fetch do corpo das mensagens.
	MailpitURL        string
	MailWebhookSecret string
}

// Load lê a configuração do ambiente, aplicando valores por omissão.
func Load() Config {
	return Config{
		HTTPAddr:        getenv("AEGIS_HTTP_ADDR", ":8080"),
		ShutdownTimeout: getenvDuration("AEGIS_SHUTDOWN_TIMEOUT", 10*time.Second),
		DatabaseURL:     getenv("AEGIS_DATABASE_URL", ""),
		AdminKey:        getenv("AEGIS_ADMIN_KEY", ""),
		MTLSAddr:        getenv("AEGIS_MTLS_ADDR", ":8443"),
		MTLSAutoDev:     getenvBool("AEGIS_MTLS_AUTO_DEV", false),
		MTLSCACert:      getenv("AEGIS_MTLS_CA_CERT", ""),
		MTLSCAKey:       getenv("AEGIS_MTLS_CA_KEY", ""),
		MTLSServerCert:  getenv("AEGIS_MTLS_SERVER_CERT", ""),
		MTLSServerKey:   getenv("AEGIS_MTLS_SERVER_KEY", ""),
		WebAuthnRPDisplayName: getenv("AEGIS_WEBAUTHN_RP_NAME", "AegisPass"),
		WebAuthnRPID:          getenv("AEGIS_WEBAUTHN_RP_ID", "localhost"),
		WebAuthnRPOrigins:     splitCSV(getenv("AEGIS_WEBAUTHN_RP_ORIGINS", "http://localhost:5173,http://localhost:8080")),
		MailpitURL:            getenv("AEGIS_MAILPIT_URL", ""),
		MailWebhookSecret:     getenv("AEGIS_MAIL_WEBHOOK_SECRET", ""),
	}
}

// getenv devolve a variável de ambiente `key` ou `fallback` se estiver vazia.
func getenv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

// getenvDuration interpreta a variável como uma duração (ex: "10s", "2m").
// Se não existir ou for inválida, usa o fallback — falhamos para o seguro.
func getenvDuration(key string, fallback time.Duration) time.Duration {
	v, ok := os.LookupEnv(key)
	if !ok || v == "" {
		return fallback
	}
	// Aceitamos tanto "10s" como um número de segundos puro ("10").
	if d, err := time.ParseDuration(v); err == nil {
		return d
	}
	if secs, err := strconv.Atoi(v); err == nil {
		return time.Duration(secs) * time.Second
	}
	return fallback
}

func getenvBool(key string, fallback bool) bool {
	v, ok := os.LookupEnv(key)
	if !ok || v == "" {
		return fallback
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return fallback
	}
	return b
}

func splitCSV(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
