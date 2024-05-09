package model

import (
	"context"
	"math"
	"math/rand"
	"time"
	"fmt"
)

// Perf controls service performance.
type Perf struct {
	Latency     *Latency
	Failure		*Failure
}

// Latency contains parameters for simulating latency
type Latency struct {
	Mean, StdDev time.Duration
}

// Failure contains parameters for simulating failures
type Failure struct {
	Probability float64
	Messages   []string
}

var defaultLatency = Latency{Mean: 15 * time.Millisecond, StdDev: 3 * time.Millisecond}

// Validate performs validation and sets defaults.
func (p *Perf) Validate(r *Registry) error {
	if p.Latency == nil {
		p.Latency = &defaultLatency
	}
	if p.Failure == nil {
		p.Failure = &Failure{Probability: 0}
	}
	return nil
}

// Apply executes the instructions specified in Perf.
func (p *Perf) Apply(context.Context) error {
	p.Latency.simulate()
	return p.Failure.simulate()
}

func (l *Latency) simulate() {
	fMean := float64(l.Mean)
	fStdDev := float64(l.StdDev)
	delay := time.Duration(math.Max(1, rand.NormFloat64()*fStdDev+fMean))
	time.Sleep(delay)
}

func (f *Failure) simulate() error {
	if rand.Float64() < f.Probability {
		return fmt.Errorf(f.Messages[rand.Intn(len(f.Messages))])
	}
	return nil
}