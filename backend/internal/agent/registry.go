package agent

import (
	"fmt"
	"sort"
	"sync"
)

// Registry mantém tools registadas por nome (thread-safe).
type Registry struct {
	mu    sync.RWMutex
	tools map[string]Tool
}

// NewRegistry cria um registo vazio.
func NewRegistry() *Registry {
	return &Registry{tools: make(map[string]Tool)}
}

// Register adiciona uma tool; falha se o nome colidir.
func (r *Registry) Register(t Tool) error {
	name := t.Descriptor().Name
	if name == "" {
		return fmt.Errorf("agent: tool sem nome")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.tools[name]; ok {
		return fmt.Errorf("%w: %s", ErrDuplicateTool, name)
	}
	r.tools[name] = t
	return nil
}

// MustRegister como Register mas faz panic em erro (arranque do servidor).
func (r *Registry) MustRegister(t Tool) {
	if err := r.Register(t); err != nil {
		panic(err)
	}
}

// Get devolve a tool pelo nome.
func (r *Registry) Get(name string) (Tool, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.tools[name]
	return t, ok
}

// ListDescriptors devolve metadados ordenados por nome (para LLM / API).
func (r *Registry) ListDescriptors() []Descriptor {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]Descriptor, 0, len(r.tools))
	for _, t := range r.tools {
		out = append(out, t.Descriptor())
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}
