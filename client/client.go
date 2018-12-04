package client

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/opentracing-contrib/go-stdlib/nethttp"
	opentracing "github.com/opentracing/opentracing-go"
)

// Client is a wrapper around http.Client that can trace the requests.
type Client struct {
	client *http.Client
}

// NewClient creates a new Client.
func NewClient() *Client {
	return &Client{
		client: &http.Client{Transport: &nethttp.Transport{}},
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

	req, ht := nethttp.TraceRequest(tracer, req)
	defer ht.Finish()

	res, err := c.client.Do(req)
	if err != nil {
		return err
	}
	res.Body.Close()
	log.Printf("status code: %d", res.StatusCode)
	if res.StatusCode != 200 {
		return fmt.Errorf("%s returned status code %d", url, res.StatusCode)
	}

	return nil
}
