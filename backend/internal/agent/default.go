package agent

// NewDefaultRegistry regista as tools de arranque (prospeção, diagnóstico).
func NewDefaultRegistry() *Registry {
	r := NewRegistry()
	r.MustRegister(NewPingTool())
	r.MustRegister(NewDraftLeadTool())
	return r
}
