package model

import "fmt"

type Registry struct {
	services map[string]*Service
}

func (r *Registry) RegisterServices(s []*Service) error {
	r.services = make(map[string]*Service)
	for _, service := range s {
		if _, ok := r.services[service.Name]; ok {
			return fmt.Errorf("Duplicate service name %s", service.Name)
		}
		r.services[service.Name] = service
	}
	return nil
}

func (r *Registry) Service(name string) *Service {
	return r.services[name]
}
