package model

// Perf controls service performance.
type Perf struct {
	FailureRate float64
}

// Validate performs validation and sets defaults.
func (p *Perf) Validate(r *Registry) error {
	return nil
}
