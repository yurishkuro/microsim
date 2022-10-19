package model_test

import (
	"testing"

	"github.com/yurishkuro/microsim/model"
)

var oneService = &model.Config{
	Services: []*model.Service{
		&model.Service{
			Name: "ui",
			Endpoints: []*model.Endpoint{
				&model.Endpoint{
					Name: "/",
				},
			},
		},
	},
}

func TestDefaults(t *testing.T) {
	cfg := oneService
	r := &model.Registry{}
	_ = r.RegisterServices(cfg.Services)

	if err := cfg.Validate(r); err != nil {
		t.Fatalf("validation failed: %v", err)
	}

	if cfg.Services[0].Endpoints[0].Perf == nil {
		t.Fatal("Perf cannot be nil")
	}
}
