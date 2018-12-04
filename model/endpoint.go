package model

import (
	"fmt"
	"net/http"
)

// Endpoint describes an endpoint
type Endpoint struct {
	Name string
	Perf *Perf `json:",omitempty"`

	service *Service
}

// Validate performs validation and sets defaults.
func (e *Endpoint) Validate(r *Registry) error {
	if e.Name == "" {
		return fmt.Errorf("service name required")
	}
	if e.Perf != nil {
		if err := e.Perf.Validate(r); err != nil {
			return fmt.Errorf("%s: perf validation error: %v", e.Name, err)
		}
	}
	return nil
}

func (e *Endpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	return
}
