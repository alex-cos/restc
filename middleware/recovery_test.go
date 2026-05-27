package middleware_test

import (
	"bytes"
	"log/slog"
	"net/http"
	"testing"

	"github.com/alex-cos/restc"
	"github.com/alex-cos/restc/middleware"
	"github.com/stretchr/testify/assert"
)

func TestRecoveryMiddlewareCatchesPanic(t *testing.T) {
	t.Parallel()

	client := restc.NewWithClient(
		"https://api.test.com",
		&mockHTTPClient{code: http.StatusOK, content: `{}`},
	)
	client.UseMiddleware(middleware.NewRecoveryMiddleware().Handler())
	client.UseMiddleware(func(req *restc.Request, next restc.HandlerFunc) (*restc.Response, error) {
		panic("something went wrong")
	})

	resp, err := client.Execute(restc.Get("users"))
	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "panic recovered")
}

func TestRecoveryMiddlewarePassThrough(t *testing.T) {
	t.Parallel()

	client := restc.NewWithClient(
		"https://api.test.com",
		&mockHTTPClient{code: http.StatusOK, content: `{"ok":true}`},
	)
	client.UseMiddleware(middleware.NewRecoveryMiddleware().Handler())
	client.UseMiddleware(func(req *restc.Request, next restc.HandlerFunc) (*restc.Response, error) {
		return next(req)
	})

	resp, err := client.Execute(restc.Get("users"))
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
}

func TestRecoveryMiddlewareLogsPanic(t *testing.T) {
	t.Parallel()

	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelError}))

	client := restc.NewWithClient(
		"https://api.test.com",
		&mockHTTPClient{code: http.StatusOK, content: `{}`},
	)
	client.UseMiddleware(middleware.NewRecoveryMiddlewareWith(logger).Handler())
	client.UseMiddleware(func(req *restc.Request, next restc.HandlerFunc) (*restc.Response, error) {
		panic("test panic")
	})

	resp, err := client.Execute(restc.Get("test"))
	assert.Error(t, err)
	assert.Nil(t, resp)

	logOutput := buf.String()
	assert.Contains(t, logOutput, "request panicked")
	assert.Contains(t, logOutput, "method=GET")
	assert.Contains(t, logOutput, "url=test")
	assert.Contains(t, logOutput, "test panic")
	assert.Contains(t, logOutput, "stack=")
}
