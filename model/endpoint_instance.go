package model

import (
	"context"
	"net/http"

	"github.com/yurishkuro/microsim/client"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/attribute"
)

// EndpointInstance implements an endpoint in a single instance of a service.
type EndpointInstance struct {
	Endpoint
	service *ServiceInstance
}

func (e *EndpointInstance) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := e.execute(r.Context())
	if err != nil {
		trace.SpanFromContext(r.Context()).SetAttributes(attribute.String("Error Message", err.Error()))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
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
func (e *EndpointInstance) Call(ctx context.Context, tracer trace.Tracer) error {
	url := e.service.server.URL + e.Name
	return client.Get(ctx, url, tracer)
}
