package restc_test

import (
	"net/http"
	"sync/atomic"
	"testing"

	"github.com/alex-cos/restc"
	"github.com/stretchr/testify/assert"
)

func TestMiddlewareSingle(t *testing.T) {
	t.Parallel()

	var called int32

	client := restc.NewWithClient("https://api.test.com",
		NewMockClient(http.StatusOK, `{"id": 1}`))

	client.UseMiddleware(func(req *restc.Request, next func(req *restc.Request) (*restc.Response, error)) (*restc.Response, error) {
		atomic.AddInt32(&called, 1)
		return next(req)
	})

	resp, err := client.Execute(restc.Get("users"))
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int32(1), atomic.LoadInt32(&called))
}

func TestMiddlewareMultipleOrder(t *testing.T) {
	t.Parallel()

	var order []int

	client := restc.NewWithClient("https://api.test.com",
		NewMockClient(http.StatusOK, `{"id": 1}`))

	client.UseMiddleware(
		func(req *restc.Request, next func(req *restc.Request) (*restc.Response, error)) (*restc.Response, error) {
			order = append(order, 1)
			resp, err := next(req)
			order = append(order, 4)
			return resp, err
		},
		func(req *restc.Request, next func(req *restc.Request) (*restc.Response, error)) (*restc.Response, error) {
			order = append(order, 2)
			resp, err := next(req)
			order = append(order, 3)
			return resp, err
		},
	)

	resp, err := client.Execute(restc.Get("users"))
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, []int{1, 2, 3, 4}, order)
}

func TestMiddlewareModifyRequest(t *testing.T) {
	t.Parallel()

	client := restc.NewWithClient("https://api.test.com",
		NewMockClient(http.StatusOK, GetResponse))

	client.UseMiddleware(func(req *restc.Request, next func(req *restc.Request) (*restc.Response, error)) (*restc.Response, error) {
		req.SetHeader("X-Middleware", "added")
		return next(req)
	})

	resp, err := client.Execute(restc.Get("users"))
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestMiddlewareModifyResponse(t *testing.T) {
	t.Parallel()

	client := restc.NewWithClient("https://api.test.com",
		NewMockClient(http.StatusOK, GetResponse))

	client.UseMiddleware(func(req *restc.Request, next func(req *restc.Request) (*restc.Response, error)) (*restc.Response, error) {
		resp, err := next(req)
		if err != nil {
			return nil, err
		}
		return resp, nil
	})

	resp, err := client.Execute(restc.Get("users"))
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 200, resp.StatusCode())
}

func TestMiddlewareShortCircuit(t *testing.T) {
	t.Parallel()

	var nextCalled bool

	client := restc.NewWithClient("https://api.test.com",
		NewMockClient(http.StatusOK, GetResponse))

	client.UseMiddleware(func(req *restc.Request, next func(req *restc.Request) (*restc.Response, error)) (*restc.Response, error) {
		return &restc.Response{}, nil
	})
	client.UseMiddleware(func(req *restc.Request, next func(req *restc.Request) (*restc.Response, error)) (*restc.Response, error) {
		nextCalled = true
		return next(req)
	})

	resp, err := client.Execute(restc.Get("users"))
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.False(t, nextCalled)
}

func TestMiddlewareError(t *testing.T) {
	t.Parallel()

	client := restc.NewWithClient("https://api.test.com",
		NewMockClient(http.StatusOK, GetResponse))

	client.UseMiddleware(func(req *restc.Request, next func(req *restc.Request) (*restc.Response, error)) (*restc.Response, error) {
		return nil, assert.AnError
	})

	resp, err := client.Execute(restc.Get("users"))
	assert.Error(t, err)
	assert.Nil(t, resp)
}
