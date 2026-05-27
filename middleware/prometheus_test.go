package middleware_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/alex-cos/restc"
	"github.com/alex-cos/restc/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
)

func TestPrometheusMiddlewareDefaultRegistry(t *testing.T) {
	t.Parallel()

	metrics := middleware.NewPrometheusMetrics("test")

	client := restc.NewWithClient(
		"https://api.test.com",
		&mockHTTPClient{code: http.StatusOK, content: `{}`},
	)
	client.UseMiddleware(metrics.Handler())

	resp, err := client.Execute(restc.Get("test"))
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
}

func TestPrometheusMiddlewareSuccess(t *testing.T) {
	t.Parallel()

	reg := prometheus.NewRegistry()
	metrics := middleware.NewPrometheusMetricsWith("test", reg)

	client := restc.NewWithClient(
		"https://api.test.com",
		&mockHTTPClient{code: http.StatusOK, content: `{"ok":true}`},
	)
	client.UseMiddleware(metrics.Handler())

	resp, err := client.Execute(restc.Get("users"))
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())

	err = testutil.CollectAndCompare(
		metrics.RequestsTotal,
		strings.NewReader(`
# HELP test_requests_total Total number of HTTP requests made.
# TYPE test_requests_total counter
test_requests_total{method="GET",status_code="200"} 1
`),
		"test_requests_total",
	)
	assert.NoError(t, err)
}

func TestPrometheusMiddlewareErrorResponse(t *testing.T) {
	t.Parallel()

	reg := prometheus.NewRegistry()
	metrics := middleware.NewPrometheusMetricsWith("test", reg)

	client := restc.NewWithClient(
		"https://api.test.com",
		&mockHTTPClient{code: http.StatusInternalServerError, content: `{"error":"fail"}`},
	)
	client.UseMiddleware(metrics.Handler())

	resp, err := client.Execute(restc.Get("users"))
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.IsError())

	err = testutil.CollectAndCompare(
		metrics.RequestsTotal,
		strings.NewReader(`
# HELP test_requests_total Total number of HTTP requests made.
# TYPE test_requests_total counter
test_requests_total{method="GET",status_code="500"} 1
`),
		"test_requests_total",
	)
	assert.NoError(t, err)
}

func TestPrometheusInFlightGauge(t *testing.T) {
	t.Parallel()

	reg := prometheus.NewRegistry()
	metrics := middleware.NewPrometheusMetricsWith("test", reg)

	client := restc.NewWithClient(
		"https://api.test.com",
		&mockHTTPClient{code: http.StatusOK, content: `{}`},
	)
	client.UseMiddleware(metrics.Handler())

	_, err := client.Execute(restc.Get("users"))
	assert.NoError(t, err)

	count := testutil.CollectAndCount(metrics.RequestsInFlight)
	assert.Equal(t, 1, count)
}
