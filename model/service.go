package model

import (
	"fmt"
	"sync/atomic"
)

// Service is service.
type Service struct {
	Name      string
	Endpoints []*Endpoint `json:",omitempty"`
	Count     int         `json:",omitempty"`

	Instances  []*ServiceInstance
	nextServer uint64
}

// Validate validates Service and sets defaults.
func (s *Service) Validate(r *Registry) error {
	if s.Name == "" {
		return fmt.Errorf("service name required: %+v", s)
	}
	if s.Count < 0 {
		return fmt.Errorf("%s: count must be > 0", s.Name)
	}
	if s.Count == 0 {
		s.Count = 1
	}
	if len(s.Endpoints) == 0 {
		s.Endpoints = []*Endpoint{defaultEndpoint()}
	}
	for i, endpoint := range s.Endpoints {
		if err := endpoint.Validate(s, r); err != nil {
			return fmt.Errorf("%s: endpoint[%d] validation error: %v", s.Name, i, err)
		}
	}
	return nil
}

// Endpoint returns an endpoint for the service, or nil if not found.
func (s *Service) Endpoint(name string) *Endpoint {
	for _, endpoint := range s.Endpoints {
		if endpoint.Name == name {
			return endpoint
		}
	}
	return nil
}

// DefaultEndpoint returns default endpoint for this service (always the first).
// If none are defined, it creates one.
func (s *Service) DefaultEndpoint() *Endpoint {
	if len(s.Endpoints) == 0 {
		s.Endpoints = []*Endpoint{defaultEndpoint()}
	}
	return s.Endpoints[0]
}

func defaultEndpoint() *Endpoint {
	return &Endpoint{
		Name: "/",
		Perf: &Perf{},
	}
}

// Start starts HTTP servers for this service.
func (s *Service) Start() error {
	s.Instances = make([]*ServiceInstance, s.Count)
	for i := 0; i < s.Count; i++ {
		instance, err := startServiceInstance(s, fmt.Sprintf("%s-%d", s.Name, i))
		if err != nil {
			return err
		}
		s.Instances[i] = instance
	}
	return nil
}

// Stop stops HTTP servers for this service.
func (s *Service) Stop() {
	for _, inst := range s.Instances {
		inst.Stop()
	}
}

// NextServerURL returns the URL of one of the servers, in round-robin fashion.
func (s *Service) NextServerURL() string {
	next := atomic.AddUint64(&s.nextServer, 1)
	nextServer := int(next) % len(s.Instances)
	return s.Instances[nextServer].server.URL
}
