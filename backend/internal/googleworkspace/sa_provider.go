package googleworkspace

import (
	"context"
	"errors"
	"os"
)

// ServiceAccountProvider valida o ficheiro JSON da SA (GOOGLE-001).
// Chamadas reais à API Google ficam para GOOGLE-002–004.
type ServiceAccountProvider struct {
	SAJSONPath    string
	DelegatedUser string
}

func (p ServiceAccountProvider) Name() string { return "service_account" }

func (p ServiceAccountProvider) Ping(_ context.Context) error {
	if p.SAJSONPath == "" || p.DelegatedUser == "" {
		return errors.New("googleworkspace: service account incompleta")
	}
	info, err := os.Stat(p.SAJSONPath)
	if err != nil {
		return err
	}
	if info.IsDir() {
		return errors.New("googleworkspace: SA_JSON é uma pasta")
	}
	// Leitura mínima — confirma que o ficheiro existe e é legível.
	f, err := os.Open(p.SAJSONPath)
	if err != nil {
		return err
	}
	defer f.Close()
	return nil
}
