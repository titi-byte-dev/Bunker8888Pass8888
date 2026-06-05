package googleworkspace

import (
	"os"
)

// Config reúne credenciais Google Workspace (GOOGLE-001).
// Didático: a Service Account é uma identidade de máquina — o JSON fica no
// servidor mas só acede a APIs; dados sensíveis continuam cifrados no cliente.
type Config struct {
	SAJSONPath    string
	DelegatedUser string
	Enabled       bool
}

// LoadConfigFromEnv lê AEGIS_GOOGLE_* — sem JSON válido usa mock provider.
func LoadConfigFromEnv() Config {
	return Config{
		SAJSONPath:    os.Getenv("AEGIS_GOOGLE_SA_JSON"),
		DelegatedUser: os.Getenv("AEGIS_GOOGLE_DELEGATED_USER"),
		Enabled:       os.Getenv("AEGIS_GOOGLE_ENABLED") == "true",
	}
}

// ProviderKind devolve "service_account" ou "mock".
func (c Config) ProviderKind() string {
	if !c.Enabled || c.SAJSONPath == "" || c.DelegatedUser == "" {
		return "mock"
	}
	if _, err := os.Stat(c.SAJSONPath); err != nil {
		return "mock"
	}
	return "service_account"
}
