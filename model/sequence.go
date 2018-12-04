package model

import "fmt"

// Sequence describes sequential dependencies.
type Sequence []Dependencies

// Validate performs validation and sets defaults.
func (s Sequence) Validate(r *Registry) error {
	for i, dep := range s {
		if err := dep.Validate(r); err != nil {
			return fmt.Errorf("Sequence[%d]: dependency validation error: %v", i, err)
		}
	}
	return nil
}
