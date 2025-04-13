package model

import (
	"context"
	"fmt"
	"log"
	"sync"

	"go.opentelemetry.io/otel/trace"
)

// Parallel describes parallel dependencies.
type Parallel struct {
	Items  []Dependencies
	MaxPar int `json:",omitempty"`
}

// Validate performs validation and sets defaults.
func (p *Parallel) Validate(r *Registry) error {
	if len(p.Items) == 0 {
		return fmt.Errorf("failed to validate Par: par requires items")
	}
	for i := range p.Items {
		if err := p.Items[i].Validate(r); err != nil {
			return fmt.Errorf("failed to validate Par: par.Item[%d]: %v", i, err)
		}
	}
	if p.MaxPar < 0 {
		return fmt.Errorf("failed to validate Par: expecting Par.MaxPar > 0")
	}
	return nil
}

// Call makes calls to all dependencies.
func (p *Parallel) Call(ctx context.Context, tracerProvider trace.TracerProvider) error {
	if p.MaxPar == 0 {
		return p.fullParCall(ctx, tracerProvider)
	}
	return p.maxParCall(ctx, tracerProvider)
}

func (p *Parallel) fullParCall(ctx context.Context, tracerProvider trace.TracerProvider) error {
	// done := &sync.WaitGroup{}
	// done.Add(len(p.Items))

	// call := func(n int) {

	// }

	// for i := range p.Items {

	// }
	log.Fatal("fullParCall not implemented") // TODO
	return nil
}

func (p *Parallel) maxParCall(ctx context.Context, tracerProvider trace.TracerProvider) error {
	done := &sync.WaitGroup{}
	done.Add(len(p.Items))

	ch := make(chan int, p.MaxPar)
	defer close(ch)

	var topErrors []error
	var errMutex sync.Mutex

	// start MaxPar goroutines
	// TODO MaxPar only affects single request, must be global
	// TODO all requests appear to end at the same time (try randomizing latency)
	for i := 0; i < p.MaxPar; i++ {
		go func() {
			for n := range ch {
				err := p.Items[n].Call(ctx, tracerProvider)
				if err != nil {
					errMutex.Lock()
					topErrors = append(topErrors, err)
					errMutex.Unlock()
				}
				done.Done()
			}
		}()
	}

	for i := range p.Items {
		ch <- i
	}

	done.Wait()

	errMutex.Lock()
	if len(topErrors) > 0 {
		return topErrors[0]
	}
	return nil
}
