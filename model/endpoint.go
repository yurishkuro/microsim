package model

import (
	"fmt"
)

// Endpoint describes an endpoint
type Endpoint struct {
	Name    string
	Perf    *Perf         `json:",omitempty"`
	Depends *Dependencies `json:",omitempty"`

	service *Service
}

// Validate performs validation and sets defaults.
func (e *Endpoint) Validate(s *Service, r *Registry) error {
	e.service = s
	if e.Name == "" {
		return fmt.Errorf("service name required")
	}
	if e.Perf != nil {
		if err := e.Perf.Validate(r); err != nil {
			return fmt.Errorf("%s: perf validation error: %v", e.Name, err)
		}
	}
	if e.Depends != nil {
		if err := e.Depends.Validate(r); err != nil {
			return fmt.Errorf("%s:%s: dependencies validation error: %v", s.Name, e.Name, err)
		}
	}
	return nil
}

// NewInstance creates a new instance of endpoint for
func (e *Endpoint) NewInstance(serviceInst *ServiceInstance) *EndpointInstance {
	return &EndpointInstance{
		Endpoint: *e,
		service:  serviceInst,
	}
}
