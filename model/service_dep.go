package model

import "fmt"

// ServiceDep describes dependency on a specific service, and optionally endpoint.
type ServiceDep struct {
	Name     string
	Endpoint string `json:",omitempty"`
}

// Validate performs validation and sets defaults.
func (s *ServiceDep) Validate(r *Registry) error {
	if s.Name == "" {
		return fmt.Errorf("Service dependency: must specify name")
	}
	service := r.Service(s.Name)
	if service == nil {
		return fmt.Errorf("Service dependency: unknown service name %s", s.Name)
	}
	if s.Endpoint == "" {
		s.Endpoint = service.DefaultEndpoint().Name
	} else {
		if service.Endpoint(s.Endpoint) == nil {
			return fmt.Errorf("Service dependency: unknown endpoint %s for service %s", s.Endpoint, s.Name)
		}
	}
	return nil
}
