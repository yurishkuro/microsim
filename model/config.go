package model

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/yurishkuro/microsim/tracing"
)

// Config is config.
type Config struct {
	Services []*Service

	TestDuration time.Duration
	TestRunners  int
	Repeats      int
}

// Validate performs validation and sets defaults.
func (c *Config) Validate(r *Registry) error {
	if len(c.Services) == 0 {
		return fmt.Errorf("Config: must have at least one service")
	}
	for i, service := range c.Services {
		if err := service.Validate(r); err != nil {
			return fmt.Errorf("Config.Service[%d]: validation error: %v", i, err)
		}
	}
	return nil
}

// Start starts the simulation.
func (c *Config) Start() error {
	for i, service := range c.Services {
		if err := service.Start(); err != nil {
			return fmt.Errorf("Start to fail service %d - %s: %v", i, service.Name, err)
		}
	}
	return nil
}

// Stop stops the simulation.
func (c *Config) Stop() {
	for _, service := range c.Services {
		service.Stop()
	}
}

// Run runs the simulation.
func (c *Config) Run() {
	stop := make(chan struct{})
	done := &sync.WaitGroup{}
	done.Add(c.TestRunners)
	for i := 0; i < c.TestRunners; i++ {
		name := fmt.Sprintf("test-executor-%d", i)
		go c.runWorker(name, stop, done)
	}
	log.Printf("started %d test executors", c.TestRunners)
	if c.Repeats > 0 {
		log.Printf("running %d repeat(s)", c.Repeats)
	} else {
		log.Printf("running for %v", c.TestDuration)
		time.Sleep(c.TestDuration)
		log.Printf("stopping test executors")
		close(stop)
	}
	log.Printf("waiting for test executors to exit")
	done.Wait()
}

func (c *Config) runWorker(instanceName string, stop chan struct{}, done *sync.WaitGroup) {
	tracer, closer, err := tracing.InitTracer("test-executor", instanceName)
	if err != nil {
		log.Fatalf("failed to create a tracer: %v", err)
	}
	defer closer.Close()
	defer done.Done()
	repeats := c.Repeats
	for {
		select {
		case <-stop:
			return
		default:
			c.runTest(tracer)
		}
		if repeats > 0 {
			repeats--
			if repeats == 0 {
				break
			}
		}
	}
}

func (c *Config) runTest(tracer opentracing.Tracer) {
	rootSvc := c.Services[0]
	inst := rootSvc.instances[0]
	endpoint := inst.Endpoints[0]
	err := endpoint.Call(context.Background(), tracer)
	if err != nil {
		log.Printf("transaction failed: %v", err)
	}
}
