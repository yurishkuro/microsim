package config

import (
	"github.com/yurishkuro/microsim/model"
)

var hotrod = &model.Config{
	Services: []*model.Service{
		{
			Name: "ui",
			Endpoints: []*model.Endpoint{
				{
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
		{
			Name: "frontend",
			Endpoints: []*model.Endpoint{
				{
					Name: "/dispatch",
					Depends: &model.Dependencies{
						Seq: model.Sequence{
							{Service: &model.ServiceDep{Name: "customer", Endpoint: "/customer"}},
							{Service: &model.ServiceDep{Name: "driver", Endpoint: "/FindNearest"}},
							{Par: &model.Parallel{
								Items: []model.Dependencies{
									{Service: &model.ServiceDep{Name: "route", Endpoint: "/GetShortestRoute"}},
									{Service: &model.ServiceDep{Name: "route", Endpoint: "/GetShortestRoute"}},
									{Service: &model.ServiceDep{Name: "route", Endpoint: "/GetShortestRoute"}},
									{Service: &model.ServiceDep{Name: "route", Endpoint: "/GetShortestRoute"}},
									{Service: &model.ServiceDep{Name: "route", Endpoint: "/GetShortestRoute"}},
									{Service: &model.ServiceDep{Name: "route", Endpoint: "/GetShortestRoute"}},
									{Service: &model.ServiceDep{Name: "route", Endpoint: "/GetShortestRoute"}},
									{Service: &model.ServiceDep{Name: "route", Endpoint: "/GetShortestRoute"}},
									{Service: &model.ServiceDep{Name: "route", Endpoint: "/GetShortestRoute"}},
									{Service: &model.ServiceDep{Name: "route", Endpoint: "/GetShortestRoute"}},
								},
								MaxPar: 3,
							}},
						},
					},
				},
			},
		},
		{
			Name: "customer",
			Endpoints: []*model.Endpoint{
				{
					Name: "/customer",
					Depends: &model.Dependencies{
						Seq: model.Sequence{
							{Service: &model.ServiceDep{Name: "mysql"}},
						},
					},
				},
			},
		},
		{
			Name: "driver",
			Endpoints: []*model.Endpoint{
				{
					Name: "/FindNearest",
					Depends: &model.Dependencies{
						Seq: model.Sequence{
							{Service: &model.ServiceDep{Name: "redis", Endpoint: "/FindDriverIDs"}},
							{Service: &model.ServiceDep{Name: "redis", Endpoint: "/GetDriver"}},
							{Service: &model.ServiceDep{Name: "redis", Endpoint: "/GetDriver"}},
							{Service: &model.ServiceDep{Name: "redis", Endpoint: "/GetDriver"}},
							{Service: &model.ServiceDep{Name: "redis", Endpoint: "/GetDriver"}},
							{Service: &model.ServiceDep{Name: "redis", Endpoint: "/GetDriver"}},
							{Service: &model.ServiceDep{Name: "redis", Endpoint: "/GetDriver"}},
							{Service: &model.ServiceDep{Name: "redis", Endpoint: "/GetDriver"}},
							{Service: &model.ServiceDep{Name: "redis", Endpoint: "/GetDriver"}},
							{Service: &model.ServiceDep{Name: "redis", Endpoint: "/GetDriver"}},
							{Service: &model.ServiceDep{Name: "redis", Endpoint: "/GetDriver"}},
						},
					},
				},
			},
			Count: 2,
		},
		{
			Name:  "route",
			Count: 3,
			Endpoints: []*model.Endpoint{
				{
					Name: "/GetShortestRoute",
				},
			},
		},
		{
			Name: "mysql",
			Endpoints: []*model.Endpoint{
				{
					Name: "/sql_select",
				},
			},
		},
		{
			Name: "redis",
			Endpoints: []*model.Endpoint{
				{
					Name: "/FindDriverIDs",
				},
				{
					Name: "/GetDriver",
					Perf: &model.Perf{
						Failure: &model.Failure{
							Probability: 0.3,
							Messages:    []string{"redis timeout", "redis error"},
						},
					},
				},
			},
		},
	},
}
