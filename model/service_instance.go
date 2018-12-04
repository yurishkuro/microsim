package model

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/opentracing-contrib/go-stdlib/nethttp"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/yurishkuro/microsim/client"
	"github.com/yurishkuro/microsim/tracing"
)

// ServiceInstance represents a single instance of a service.
type ServiceInstance struct {
	Endpoints []*EndpointInstance
	service   *Service
	server    *httptest.Server
	client    *client.Client
	tracing   struct {
		tracer opentracing.Tracer
		closer io.Closer
	}
}

func startServiceInstance(service *Service, instanceName string) (*ServiceInstance, error) {
	tracer, closer, err := tracing.InitTracer(instanceName)
	if err != nil {
		return nil, err
	}
	inst := &ServiceInstance{
		service: service,
		client:  client.NewClient(),
	}
	inst.tracing.tracer = tracer
	inst.tracing.closer = closer
	inst.server = httptest.NewServer(inst.mux())
	log.Printf("started service instance %s at %s", instanceName, inst.server.URL)
	return inst, nil
}

func (inst *ServiceInstance) mux() http.Handler {
	mux := http.NewServeMux()
	inst.Endpoints = make([]*EndpointInstance, len(inst.service.Endpoints))
	for i, endpoint := range inst.service.Endpoints {
		inst.Endpoints[i] = endpoint.NewInstance(inst)
		mw := nethttp.Middleware(
			inst.tracing.tracer,
			inst.Endpoints[i],
			nethttp.OperationNameFunc(func(r *http.Request) string {
				return endpoint.Name
			}))
		mux.Handle(endpoint.Name, mw)
	}
	return mux
}

// Stop shuts down the HTTP server and closes the tracer.
func (inst *ServiceInstance) Stop() {
	inst.server.Close()
	inst.tracing.closer.Close()
	log.Printf("stopped service instance %s", inst.service.Name)
}
