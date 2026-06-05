package googleworkspace

import "context"

// MockProvider simula Workspace sem chamar APIs Google (dev / sem VPS).
type MockProvider struct{}

func (MockProvider) Name() string { return "mock" }

func (MockProvider) Ping(_ context.Context) error { return nil }
