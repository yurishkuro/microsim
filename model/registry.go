package model

import "fmt"

// Registry stores all services.
type Registry struct {
	services map[string]*Service
}

// RegisterServices stores all services in Registry.
func (r *Registry) RegisterServices(s []*Service) error {
	r.services = make(map[string]*Service)
	for _, service := range s {
		if _, ok := r.services[service.Name]; ok {
			return fmt.Errorf("duplicate service name %s", service.Name)
		}
		r.services[service.Name] = service
	}
	return nil
}

// Service looks up a service by name.
func (r *Registry) Service(name string) *Service {
	return r.services[name]
}
