package tracing

import (
	"io"

	opentracing "github.com/opentracing/opentracing-go"
	jaeger "github.com/uber/jaeger-client-go"
	jaegerCfg "github.com/uber/jaeger-client-go/config"
)

// JaegerCollectorURL defines the address of Jaeger collector to submit spans.
var JaegerCollectorURL = "http://localhost:14268/api/traces"

// InitTracer creates a new tracer for a service.
func InitTracer(serviceName, instanceName string) (opentracing.Tracer, io.Closer, error) {
	cfg := &jaegerCfg.Configuration{
		ServiceName: serviceName,
		Sampler: &jaegerCfg.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &jaegerCfg.ReporterConfig{
			// LogSpans: true,
			CollectorEndpoint: JaegerCollectorURL,
		},
		Tags: []opentracing.Tag{
			{Key: "service-instance", Value: instanceName},
		},
	}
	return cfg.NewTracer(jaegerCfg.Logger(jaeger.StdLogger))
}
