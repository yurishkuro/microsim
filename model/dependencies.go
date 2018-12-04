package model

import (
	"fmt"
)

// Dependencies describes dependencies.
type Dependencies struct {
	Seq     Sequence    `json:",omitempty"`
	Par     *Parallel   `json:",omitempty"`
	Service *ServiceDep `json:",omitempty"`
}

// Validate performs validation and sets defaults.
func (d *Dependencies) Validate(r *Registry) error {
	count := 0
	if len(d.Seq) > 0 {
		if err := d.Seq.Validate(r); err != nil {
			return fmt.Errorf("Dependencies.Seq validation error: %v", err)
		}
		count++
	}
	if d.Par != nil {
		if err := d.Par.Validate(r); err != nil {
			return fmt.Errorf("Dependencies.Par validation error: %v", err)
		}
		count++
	}
	if d.Service != nil {
		if err := d.Service.Validate(r); err != nil {
			return fmt.Errorf("Dependencies.Service validation error: %v", err)
		}
		count++
	}
	if count != 1 {
		return fmt.Errorf("Dependencies: exactly one of the fields must be populated")
	}
	return nil
}
