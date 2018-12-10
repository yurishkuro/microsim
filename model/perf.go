package model

import (
	"context"
	"math"
	"math/rand"
	"time"
)

// Perf controls service performance.
type Perf struct {
	Latency     *Latency
	FailureRate float64
}

// Latency contains parameters for simulating latency
type Latency struct {
	Mean, StdDev time.Duration
}

var defaultLatency = Latency{Mean: 15 * time.Millisecond, StdDev: 3 * time.Millisecond}

// Validate performs validation and sets defaults.
func (p *Perf) Validate(r *Registry) error {
	if p.Latency == nil {
		p.Latency = &defaultLatency
	}
	return nil
}

// Apply executes the instructions specified in Perf.
func (p *Perf) Apply(context.Context) error {
	p.Latency.simulate()
	// TODO implement failures
	return nil
}

func (l *Latency) simulate() {
	fMean := float64(l.Mean)
	fStdDev := float64(l.StdDev)
	delay := time.Duration(math.Max(1, rand.NormFloat64()*fStdDev+fMean))
	time.Sleep(delay)
}
