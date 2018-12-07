package model

import (
	"context"
	"fmt"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/yurishkuro/microsim/client"
)

// ServiceDep describes dependency on a specific service, and optionally endpoint.
type ServiceDep struct {
	Name     string
	Endpoint string `json:",omitempty"`

	service  *Service
	endpoint *Endpoint
}

// Validate performs validation and sets defaults.
func (s *ServiceDep) Validate(r *Registry) error {
	if s.Name == "" {
		return fmt.Errorf("Service dependency: must specify name")
	}
	s.service = r.Service(s.Name)
	if s.service == nil {
		return fmt.Errorf("Service dependency: unknown service name %s", s.Name)
	}
	if s.Endpoint == "" {
		s.Endpoint = s.service.DefaultEndpoint().Name
	}
	s.endpoint = s.service.Endpoint(s.Endpoint)
	if s.endpoint == nil {
		return fmt.Errorf("Service dependency: unknown endpoint %s for service %s", s.Endpoint, s.Name)
	}
	return nil
}

// Call makes call to dependency service.
func (s *ServiceDep) Call(ctx context.Context, tracer opentracing.Tracer) error {
	url := s.service.NextServerURL() + s.endpoint.Name
	return client.Get(ctx, url, tracer)
}
