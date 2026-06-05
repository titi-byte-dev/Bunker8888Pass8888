package googleworkspace

import "context"

// Status resume a ligação Google Workspace para a UI.
type Status struct {
	Provider      string   `json:"provider"`
	Enabled       bool     `json:"enabled"`
	DelegatedUser string   `json:"delegated_user,omitempty"`
	Ready         bool     `json:"ready"`
	Scopes        []string `json:"scopes"`
	Message       string   `json:"message,omitempty"`
}

// Service expõe estado do provider activo (mock ou service account).
type Service struct {
	Config   Config
	Provider Provider
}

func NewService(cfg Config, provider Provider) *Service {
	return &Service{Config: cfg, Provider: provider}
}

// Status verifica se o provider responde (Ping) e devolve metadados seguros.
func (s *Service) Status(ctx context.Context) Status {
	if s == nil || s.Provider == nil {
		return Status{Provider: "disabled", Message: "google workspace indisponível"}
	}
	st := Status{
		Provider:      s.Provider.Name(),
		Enabled:       s.Config.Enabled,
		DelegatedUser: s.Config.DelegatedUser,
		Scopes: []string{
			"https://www.googleapis.com/auth/drive",
			"https://www.googleapis.com/auth/spreadsheets",
			"https://www.googleapis.com/auth/gmail.send",
		},
	}
	if err := s.Provider.Ping(ctx); err != nil {
		st.Ready = false
		st.Message = err.Error()
		return st
	}
	st.Ready = true
	if st.Provider == "mock" {
		st.Message = "modo simulação — usa /work/google-dev até OAuth na VPS"
	}
	return st
}

// SelectProvider escolhe mock ou service account conforme env.
func SelectProvider(cfg Config) Provider {
	if cfg.ProviderKind() == "service_account" {
		return ServiceAccountProvider{
			SAJSONPath:    cfg.SAJSONPath,
			DelegatedUser: cfg.DelegatedUser,
		}
	}
	return MockProvider{}
}
