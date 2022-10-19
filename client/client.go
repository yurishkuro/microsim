package client

import (
	"context"
	"fmt"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/trace"
)

// Get makes a traced HTTP GET call.
func Get(ctx context.Context, url string, tracer trace.Tracer) error {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)

	if err != nil {
		return err
	}

	client := http.Client{
		Transport: otelhttp.NewTransport(
			http.DefaultTransport,
			otelhttp.WithTracerProvider(&tracerProvider{tracer}),
		),
	}
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	res.Body.Close()
	if res.StatusCode != 200 {
		return fmt.Errorf("%s returned status code %d", url, res.StatusCode)
	}

	return nil
}

type tracerProvider struct {
	tracer trace.Tracer
}

func (p *tracerProvider) Tracer(_ string, _ ...trace.TracerOption) trace.Tracer {
	return p.tracer
}
