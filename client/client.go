package client

import (
	"context"
	"fmt"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/trace"
)

// Get makes a traced HTTP GET call.
func Get(ctx context.Context, url string, tp trace.TracerProvider) error {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}

	client := http.Client{
		Transport: otelhttp.NewTransport(
			http.DefaultTransport,
			otelhttp.WithTracerProvider(tp),
		),
	}
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	if err := res.Body.Close(); err != nil {
		return fmt.Errorf("failed to close response body: %w", err)
	}
	if res.StatusCode != 200 {
		return fmt.Errorf("%s returned status code %d", url, res.StatusCode)
	}

	return nil
}
