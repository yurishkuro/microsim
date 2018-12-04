package model

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
)

// Service is service.
type Service struct {
	Name      string
	Endpoints []*Endpoint   `json:",omitempty"`
	Count     int           `json:",omitempty"`
	Depends   *Dependencies `json:",omitempty"`

	servers []*httptest.Server
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
		if err := endpoint.Validate(r); err != nil {
			return fmt.Errorf("%s: endpoint[%d] validation error: %v", s.Name, i, err)
		}
		endpoint.service = s
	}
	if s.Depends != nil {
		if err := s.Depends.Validate(r); err != nil {
			return fmt.Errorf("%s: dependencies validation error: %v", s.Name, err)
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

func defaultEndpoint() *Endpoint {
	return &Endpoint{
		Name: "/",
		Perf: &Perf{},
	}
}

// Start starts HTTP servers for this service.
func (s *Service) Start() error {
	s.servers = make([]*httptest.Server, s.Count)
	for i := 0; i < s.Count; i++ {
		s.servers[i] = httptest.NewServer(s.mux())
		log.Printf("started service %s[%d] at %s", s.Name, i+1, s.servers[i].URL)
	}
	return nil
}

func (s *Service) mux() http.Handler {
	mux := http.NewServeMux()
	for _, endpoint := range s.Endpoints {
		mux.Handle(endpoint.Name, endpoint)
	}
	return mux
}
