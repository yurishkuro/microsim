package model

import (
	"context"
	"fmt"

	opentracing "github.com/opentracing/opentracing-go"
)

// Sequence describes sequential dependencies.
type Sequence []Dependencies

// Validate performs validation and sets defaults.
func (s Sequence) Validate(r *Registry) error {
	for i := range s {
		if err := s[i].Validate(r); err != nil {
			return fmt.Errorf("Sequence[%d]: dependency validation error: %v", i, err)
		}
	}
	return nil
}

// Call makes calls to all dependencies.
func (s Sequence) Call(ctx context.Context, tracer opentracing.Tracer) error {
	for _, dep := range s {
		if err := dep.Call(ctx, tracer); err != nil {
			return err
		}
	}
	return nil
}
