package client

import (
	"context"
	"fmt"
	"log"
	"net/http"

	opentracing "github.com/opentracing/opentracing-go"
	ottag "github.com/opentracing/opentracing-go/ext"
	otlog "github.com/opentracing/opentracing-go/log"
)

// Client is a wrapper around http.Client that can trace the requests.
type Client struct {
	client *http.Client
}

// NewClient creates a new Client.
func NewClient() *Client {
	return &Client{
		client: &http.Client{},
	}
}

// Call makes a traced HTTP call.
func (c *Client) Call(ctx context.Context, url string, tracer opentracing.Tracer) error {
	log.Printf("calling url: %s", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req = req.WithContext(ctx)

	parent := opentracing.SpanFromContext(req.Context())
	var parentCtx opentracing.SpanContext
	if parent != nil {
		parentCtx = parent.Context()
	}
	span := tracer.StartSpan("HTTP GET", opentracing.ChildOf(parentCtx))
	ottag.SpanKindRPCClient.Set(span)
	ottag.HTTPMethod.Set(span, req.Method)
	ottag.HTTPUrl.Set(span, req.URL.String())
	defer span.Finish()

	carrier := opentracing.HTTPHeadersCarrier(req.Header)
	tracer.Inject(span.Context(), opentracing.HTTPHeaders, carrier)

	res, err := c.client.Do(req)
	if err != nil {
		ottag.Error.Set(span, true)
		span.LogFields(otlog.Error(err))
		return err
	}

	res.Body.Close()
	ottag.HTTPStatusCode.Set(span, uint16(res.StatusCode))
	log.Printf("status code: %d", res.StatusCode)
	if res.StatusCode != 200 {
		return fmt.Errorf("%s returned status code %d", url, res.StatusCode)
	}

	return nil
}
