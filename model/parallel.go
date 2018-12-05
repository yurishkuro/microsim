package model

import (
	"context"
	"log"

	opentracing "github.com/opentracing/opentracing-go"
)

// Parallel describes parallel dependencies.
type Parallel struct {
	Seq     Sequence    `json:",omitempty"`
	Service *ServiceDep `json:",omitempty"`
	MaxPar  int         `json:",omitempty"`
}

// Validate performs validation and sets defaults.
func (p *Parallel) Validate(r *Registry) error {
	log.Fatal("not implemented")
	return nil
}

// Call makes calls to all dependencies.
func (p *Parallel) Call(ctx context.Context, tracer opentracing.Tracer) error {
	log.Fatal("not implemented")
	return nil
}
