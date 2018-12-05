package model

import (
	"context"
	"net/http"

	opentracing "github.com/opentracing/opentracing-go"
)

// EndpointInstance implements an endpoint in a single instance of a service.
type EndpointInstance struct {
	Endpoint
	service *ServiceInstance
}

func (e *EndpointInstance) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := e.execute(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - " + err.Error()))
	}
	return
}

// executes the endpoint, calling dependencies if necessary.
func (e *EndpointInstance) execute(ctx context.Context) error {
	if e.Depends != nil {
		if err := e.Depends.Call(ctx, e.service.tracing.tracer); err != nil {
			return err
		}
	}
	return e.Perf.Apply(ctx)
}

// Call makes a call to this endpoint.
func (e *EndpointInstance) Call(ctx context.Context, tracer opentracing.Tracer) error {
	url := e.service.server.URL + e.Name
	return e.service.client.Call(ctx, url, tracer)
}
