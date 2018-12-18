package client

import (
	"context"
	"fmt"
	"net/http"

	opentracing "github.com/opentracing/opentracing-go"
	ottag "github.com/opentracing/opentracing-go/ext"
	otlog "github.com/opentracing/opentracing-go/log"
)

// Get makes a traced HTTP GET call.
func Get(ctx context.Context, url string, tracer opentracing.Tracer) error {
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

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		ottag.Error.Set(span, true)
		span.LogFields(otlog.Error(err))
		return err
	}

	res.Body.Close()
	ottag.HTTPStatusCode.Set(span, uint16(res.StatusCode))
	if res.StatusCode != 200 {
		return fmt.Errorf("%s returned status code %d", url, res.StatusCode)
	}

	return nil
}
