package main

import (
	"github.com/yurishkuro/microsim/model"
)

var hotrod = model.Config{
	Services: []*model.Service{
		&model.Service{
			Name: "ui",
			Endpoints: []*model.Endpoint{
				&model.Endpoint{
					Name: "/",
					Depends: &model.Dependencies{
						Service: &model.ServiceDep{
							Name:     "frontend",
							Endpoint: "/dispatch",
						},
					},
				},
			},
		},
		&model.Service{
			Name: "frontend",
			Endpoints: []*model.Endpoint{
				&model.Endpoint{
					Name: "/dispatch",
					Depends: &model.Dependencies{
						Seq: model.Sequence{
							{Service: &model.ServiceDep{Name: "customer"}},
							{Service: &model.ServiceDep{Name: "driver"}},
							{Service: &model.ServiceDep{Name: "route"}},
						},
					},
				},
			},
		},
		&model.Service{
			Name: "customer",
			Endpoints: []*model.Endpoint{
				&model.Endpoint{
					Name: "/sql_select",
					Depends: &model.Dependencies{
						Seq: model.Sequence{
							{Service: &model.ServiceDep{Name: "mysql"}},
						},
					},
				},
			},
		},
		&model.Service{
			Name:  "driver",
			Count: 2,
		},
		&model.Service{
			Name:  "route",
			Count: 3,
		},
		&model.Service{
			Name: "mysql",
		},
		&model.Service{
			Name: "redis",
		},
	},
}
