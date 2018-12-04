package model

import "context"

// Perf controls service performance.
type Perf struct {
	FailureRate float64
}

// Validate performs validation and sets defaults.
func (p *Perf) Validate(r *Registry) error {
	return nil
}

// Apply executes the instructions specified in Perf.
func (p *Perf) Apply(context.Context) error {
	return nil
}
