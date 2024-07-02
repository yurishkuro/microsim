package tracing

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.opentelemetry.io/otel/attribute"
	otel "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	"go.opentelemetry.io/otel/trace"
)

// InitTracer creates a new tracer for a service.
func InitTracer(serviceName, instanceName string) (trace.TracerProvider, func(), error) {
	tp, err := newTracerProvider(serviceName, instanceName)
	if err != nil {
		return nil, nil, fmt.Errorf("couldn't initialize tracer provider: %w", err)
	}

	// Cleanly shutdown and flush telemetry when the application exits.
	shutdown := func() {
		// Do not make the application hang when it is shutdown.
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		if err := tp.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}
	}
	return tp, shutdown, err
}

// tracerProvider returns an OpenTelemetry TracerProvider configured to use
// the OTEL exporter that will send spans to the provided url. The returned
// TracerProvider will also use a Resource configured with all the information
// about the application.
func newTracerProvider(serviceName, instanceName string) (*tracesdk.TracerProvider, error) {
	// Create the OTEL exporter
	exp, err := otel.New(context.Background(), otel.WithEndpointURL(os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")))
	if err != nil {
		return nil, err
	}
	return tracesdk.NewTracerProvider(
		tracesdk.WithSampler(tracesdk.AlwaysSample()),
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),
		// Record information about this application in a Resource.
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
			attribute.String("service-instance", instanceName),
		)),
	), nil
}
