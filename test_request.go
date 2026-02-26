package scuter

import (
	"bytes"
	"context"
	"encoding/json/v2"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
)

// NewTestRequest returns a new incoming server *http.Request, suitable for passing to an [http.Handler] for testing.
// IMPORTANT: Any error encountered in building the *http.Request results in a panic.
func NewTestRequest(ctx context.Context, method, target string, options ...RequestOption) *http.Request {
	config := requestConfig{
		body:    bytes.NewBuffer(nil),
		headers: make(http.Header),
		query:   make(url.Values),
	}
	Request.With(options...)(&config)

	request := httptest.NewRequestWithContext(ctx, method, target, config.body)

	query := request.URL.Query()
	for key, values := range config.query {
		for _, value := range values {
			query.Add(key, value)
		}
	}
	request.URL.RawQuery = query.Encode()

	for key, values := range config.headers {
		for _, value := range values {
			request.Header.Add(key, value)
		}
	}

	return request
}

// RequestOption is a callback func with an opportunity to modify the *requestConfig.
type RequestOption func(*requestConfig)

// Request is the 'namespace' for all methods that return a RequestOption.
var Request requestSingleton

type requestSingleton struct{}

// With returns a 'composite' option which will be the result of calling all options in the provided order.
func (requestSingleton) With(options ...RequestOption) RequestOption {
	return func(config *requestConfig) {
		for _, option := range options {
			if option != nil {
				option(config)
			}
		}
	}
}

// Query returns an option which will add the provided key/value to the request URL query string.
func (requestSingleton) Query(key, value string) RequestOption {
	return func(c *requestConfig) { c.query.Add(key, value) }
}

// Header returns an option which will add the provided key/value to the request header.
func (requestSingleton) Header(key, value string) RequestOption {
	return func(c *requestConfig) { c.headers.Add(key, value) }
}

// Body returns an option which will copy the provided reader to the request body.
func (requestSingleton) Body(r io.Reader) RequestOption {
	return func(c *requestConfig) {
		_, err := io.Copy(c.body, r)
		if err != nil {
			panic(err)
		}
	}
}

// JSONBody returns an option which will marshal the provided value to the request body with a JSON content-type header.
func (requestSingleton) JSONBody(v any) RequestOption {
	return func(c *requestConfig) {
		Request.Header(headerContentType, jsonContentType)(c)
		err := json.MarshalWrite(c.body, v)
		if err != nil {
			panic(err)
		}
	}
}

type requestConfig struct {
	query   url.Values
	headers http.Header
	body    *bytes.Buffer
}
