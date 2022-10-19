package model

import (
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/yurishkuro/microsim/tracing"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/trace"
)

// ServiceInstance represents a single instance of a service.
type ServiceInstance struct {
	Endpoints []*EndpointInstance
	service   *Service
	server    *httptest.Server
	tracing   struct {
		tracer   trace.Tracer
		shutdown func()
	}
}

func startServiceInstance(service *Service, instanceName string) (*ServiceInstance, error) {
	tracer, shutdown, err := tracing.InitTracer(service.Name, instanceName)
	if err != nil {
		return nil, err
	}
	inst := &ServiceInstance{
		service: service,
	}
	inst.tracing.tracer = tracer
	inst.tracing.shutdown = shutdown
	inst.server = httptest.NewServer(inst.mux())
	log.Printf("started service instance %s at %s", instanceName, inst.server.URL)
	return inst, nil
}

func (inst *ServiceInstance) mux() http.Handler {
	mux := http.NewServeMux()
	inst.Endpoints = make([]*EndpointInstance, len(inst.service.Endpoints))
	for i, endpoint := range inst.service.Endpoints {
		endpointInstance := endpoint.NewInstance(inst)
		inst.Endpoints[i] = endpointInstance
		wrappedHandler := otelhttp.NewHandler(
			endpointInstance,
			endpointInstance.Name,
			otelhttp.WithTracerProvider(newSingletonTracerProvider(inst.tracing.tracer)),
			otelhttp.WithSpanNameFormatter(func(_ string, _ *http.Request) string {
				return endpointInstance.Name
			}),
		)
		mux.Handle(endpoint.Name, wrappedHandler)
	}
	return mux
}

// Stop shuts down the HTTP server and closes the tracer.
func (inst *ServiceInstance) Stop() {
	inst.server.Close()
	inst.tracing.shutdown()
	log.Printf("stopped service instance %s", inst.service.Name)
}

type singletonTracerProvider struct {
	tracer trace.Tracer
}

func (p *singletonTracerProvider) Tracer(_ string, _ ...trace.TracerOption) trace.Tracer {
	return p.tracer
}

func newSingletonTracerProvider(tracer trace.Tracer) trace.TracerProvider {
	return &singletonTracerProvider{tracer}
}
