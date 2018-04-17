package health

import (
	"net/http"
	"time"

	"github.com/pkg/errors"
)

// HTTPHealthCheckOptions defines options for the http health check
type HTTPHealthCheckOptions struct {
	URL                string
	Method             string
	Username           string
	Password           string
	Client             *http.Client
	Timeout            time.Duration
	Header             http.Header
	ExpectedStatusCode int
}

// HTTPHealthCheckBuilder
type HTTPHealthCheckBuilder struct {
	options *HTTPHealthCheckOptions
}

// WithBasicAuth adds basic authentication to the test request
func (builder *HTTPHealthCheckBuilder) WithBasicAuth(username string, password string) *HTTPHealthCheckBuilder {
	builder.options.Username = username
	builder.options.Password = password
	return builder
}

// WithMethod sets the method for the test request
func (builder *HTTPHealthCheckBuilder) WithMethod(method string) *HTTPHealthCheckBuilder {
	builder.options.Method = method
	return builder
}

// WithExpectedStatusCode sets the expected status code for the test request
func (builder *HTTPHealthCheckBuilder) WithExpectedStatusCode(expectedStatusCode int) *HTTPHealthCheckBuilder {
	builder.options.ExpectedStatusCode = expectedStatusCode
	return builder
}

// WithHttpClient sets the http client which is used to execute the test request
func (builder *HTTPHealthCheckBuilder) WithHttpClient(client *http.Client) *HTTPHealthCheckBuilder {
	builder.options.Client = client
	return builder
}

// WithTimeout sets the timeout for the http request. Note this method does not work with the WithHttpClient method.
func (builder *HTTPHealthCheckBuilder) WithTimeout(timeout time.Duration) *HTTPHealthCheckBuilder {
	builder.options.Timeout = timeout
	return builder
}

// WithHeader sets an header for the test request
func (builder *HTTPHealthCheckBuilder) WithHeader(key string, value string) *HTTPHealthCheckBuilder {
	builder.options.Header.Add(key, value)
	return builder
}

// Builder creates a new http health checker with the configured options
func (builder *HTTPHealthCheckBuilder) Build() HealthChecker {
	return NewHTTPHealthCheckerWithOptions(builder.options)
}

// NewHTTPHealthCheckBuilder creates new options for the http health check with useful defaults
func NewHTTPHealthCheckBuilder(url string) *HTTPHealthCheckBuilder {
	return &HTTPHealthCheckBuilder{
		options: createHTTPHealthCheckOptions(url),
	}
}

// NewHTTPHealthChecker creates a new http health check with useful defaults
func NewHTTPHealthChecker(url string) HealthChecker {
	return NewHTTPHealthCheckerWithOptions(createHTTPHealthCheckOptions(url))
}

func createHTTPHealthCheckOptions(url string) *HTTPHealthCheckOptions {
	return &HTTPHealthCheckOptions{
		URL:                url,
		Method:             "HEAD",
		Timeout:            30 * time.Second,
		ExpectedStatusCode: 200,
		Header:             http.Header{},
	}
}

// NewHTTPHealthCheckerWithOptions creates a new http health check with the given options
func NewHTTPHealthCheckerWithOptions(options *HTTPHealthCheckOptions) HealthChecker {
	return func() error {
		client := createHttpClient(options)

		req, err := createHttpRequest(options)
		if err != nil {
			return err
		}

		resp, err := client.Do(req)
		if err != nil {
			return errors.Wrapf(err, "request for %s failed", options.URL)
		}

		if resp.StatusCode != options.ExpectedStatusCode {
			return errors.Errorf("request for %s returned %v instead of expected %v", options.URL, resp.StatusCode, options.ExpectedStatusCode)
		}

		return nil
	}
}

func createHttpRequest(options *HTTPHealthCheckOptions) (*http.Request, error) {
	req, err := http.NewRequest(options.Method, options.URL, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create request for %s: %s", options.Method, options.URL)
	}

	if options.Username != "" {
		req.SetBasicAuth(options.Username, options.Password)
	}

	for key, values := range options.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}
	return req, err
}

func createHttpClient(options *HTTPHealthCheckOptions) *http.Client {
	client := options.Client
	if client == nil {
		client = &http.Client{
			Timeout: options.Timeout,
		}
	}
	return client
}
