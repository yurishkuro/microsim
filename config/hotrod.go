package config

import (
	"github.com/yurishkuro/microsim/model"
)

var hotrod = &model.Config{
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
		&model.Service{
			Name: "customer",
			Endpoints: []*model.Endpoint{
				&model.Endpoint{
					Name: "/customer",
					Depends: &model.Dependencies{
						Seq: model.Sequence{
							{Service: &model.ServiceDep{Name: "mysql"}},
						},
					},
				},
			},
		},
		&model.Service{
			Name: "driver",
			Endpoints: []*model.Endpoint{
				&model.Endpoint{
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
		&model.Service{
			Name:  "route",
			Count: 3,
			Endpoints: []*model.Endpoint{
				&model.Endpoint{
					Name: "/GetShortestRoute",
				},
			},
		},
		&model.Service{
			Name: "mysql",
			Endpoints: []*model.Endpoint{
				&model.Endpoint{
					Name: "/sql_select",
				},
			},
		},
		&model.Service{
			Name: "redis",
			Endpoints: []*model.Endpoint{
				&model.Endpoint{
					Name: "/FindDriverIDs",
				},
				&model.Endpoint{
					Name: "/GetDriver",
				},
			},
		},
	},
}
