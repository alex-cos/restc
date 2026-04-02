package middleware_test

import (
	"bytes"
	"context"
	"log/slog"
	"net/http"
	"strings"
	"testing"

	"github.com/alex-cos/restc"
	"github.com/alex-cos/restc/middleware"
	"github.com/stretchr/testify/assert"
)

type mockHTTPClient struct {
	code    int
	content string
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return &http.Response{
		Status:     "200 OK",
		StatusCode: m.code,
		Proto:      "HTTP/2.0",
		Header: http.Header{
			"Content-Type": []string{"application/json"},
		},
		Body: &readCloser{strings.NewReader(m.content)},
	}, nil
}

type readCloser struct {
	*strings.Reader
}

func (r *readCloser) Close() error { return nil }

func TestLoggerMiddlewareSuccess(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelInfo}))

	client := restc.NewWithClient(
		"https://api.test.com",
		&mockHTTPClient{code: http.StatusOK, content: `{"ok":true}`},
	)
	client.UseMiddleware(middleware.NewMiddlewareLoggerWith(logger).Handler())

	resp, err := client.Execute(restc.Get("users"))
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())

	logOutput := buf.String()
	assert.Contains(t, logOutput, "method=GET")
	assert.Contains(t, logOutput, "url=users")
	assert.Contains(t, logOutput, "status=200")
	assert.Contains(t, logOutput, "HTTP request completed")
}

func TestLoggerMiddlewareError(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug}))

	client := restc.NewWithClient(
		"https://api.test.com",
		&mockHTTPClient{code: http.StatusInternalServerError, content: `{"error":"fail"}`},
	)
	client.UseMiddleware(middleware.NewMiddlewareLoggerWith(logger).Handler())

	resp, err := client.Execute(restc.Get("users"))
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.IsError())

	logOutput := buf.String()
	assert.Contains(t, logOutput, "status=500")
	assert.Contains(t, logOutput, "HTTP error response")
}

func TestLoggerMiddlewareWithContext(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelInfo}))

	client := restc.NewWithClient(
		"https://api.test.com",
		&mockHTTPClient{code: http.StatusOK, content: `{}`},
	)
	client.UseMiddleware(middleware.NewMiddlewareLoggerWith(logger).Handler())

	ctx := context.Background()
	resp, err := client.ExecuteWithContext(ctx, restc.Get("users"))
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	logOutput := buf.String()
	assert.Contains(t, logOutput, "HTTP request completed")
}
