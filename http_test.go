package health_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cloudogu/go-health"
	"github.com/stretchr/testify/assert"
)

func TestNewHTTPHealthChecker(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/health" || r.Method != "HEAD" {
			w.WriteHeader(404)
			return
		}

		w.WriteHeader(200)
	}))
	defer server.Close()

	checker := health.NewHTTPHealthChecker(server.URL + "/health")
	err := checker()
	assert.NoError(t, err)
}

func TestNewHTTPHealthCheckerWithOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/healthy" {
			w.WriteHeader(404)
			return
		}

		if r.Method != "POST" {
			w.WriteHeader(405)
			return
		}

		if r.Header.Get("Content-Type") != "nothing" {
			w.WriteHeader(406)
			return
		}

		if username, password, _ := r.BasicAuth(); username != "trillian" || password != "tricia123" {
			w.WriteHeader(401)
			return
		}

		w.WriteHeader(204)
	}))
	defer server.Close()

	checker := health.NewHTTPHealthCheckBuilder(server.URL+"/healthy").
		WithMethod("POST").
		WithHeader("Content-Type", "nothing").
		WithBasicAuth("trillian", "tricia123").
		WithExpectedStatusCode(204).
		Build()

	err := checker()
	assert.NoError(t, err)
}

func TestNewHTTPHealthCheckerWithOptions_ExpectedStatusCode(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(204)
	}))
	defer server.Close()
	checker := health.NewHTTPHealthCheckBuilder(server.URL + "/healthy").
		WithExpectedStatusCode(200).
		Build()

	err := checker()
	assert.Error(t, err)
}

func TestNewHTTPHealthChecker_WithFailedRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))
	checker := health.NewHTTPHealthChecker(server.URL + "/healthy")
	server.Close()
	err := checker()
	assert.Error(t, err)
}
