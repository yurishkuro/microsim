package model

import (
	"fmt"
)

// Config is config.
type Config struct {
	Services []*Service
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
