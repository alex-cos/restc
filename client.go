// Package restc provides a simple and flexible HTTP client for making REST API requests.
//
// # Features
//
// - Fluent API for building requests with method chaining
// - Automatic retry with exponential backoff
// - Configurable redirect handling
// - Middleware support for request/response interception
// - Support for JSON, XML, and multipart form data
// - IPv4/IPv6 transport options
//

package restc

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

const (
	// DefaultTimeout is the default timeout for HTTP requests.
	DefaultTimeout = 12 * time.Second
	// DefaultWaitTime is the default wait time between retries.
	DefaultWaitTime = time.Duration(100) * time.Millisecond
	// DefaultMaxWaitTime is the maximum wait time between retries.
	DefaultMaxWaitTime = time.Duration(2000) * time.Millisecond
	// DefaultUserAgent is the default "User-Agent" value.
	DefaultUserAgent = "restc/1.0"
)

// HTTPClient interface represents an HTTP client capable of executing requests.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// ParseResponse is a function that parses an HTTP response into a custom type.
type ParseResponse func(request *Request, response *Response) (any, error)

// Client is an HTTP client for making REST API requests.
// It provides methods for configuring and executing HTTP requests.
type Client struct {
	entryPoint       string
	client           HTTPClient
	timeout          time.Duration
	retryCount       int
	retryWaitTime    time.Duration
	retryMaxWaitTime time.Duration
	maxResponseSize  int64
	parseResponse    ParseResponse
	parseError       ParseResponse
	middleware       *ClientMiddleware
	redirectConfig   RedirectConfig
	defaultHeaders   map[string]string
	mutex            *sync.RWMutex
}

// New creates a new Client with the given entry point URL and options.
// The entry point is used as the base URL for relative paths in requests.
func New(entryPoint string, opts ...Option) *Client {
	return NewWithClient(entryPoint, http.DefaultClient, opts...)
}

// NewWithClient creates a new Client with the given entry point URL,
// custom HTTP client, and options.
// This allows using a custom http.Client configuration.
func NewWithClient(entryPoint string, httpClient HTTPClient, opts ...Option) *Client {
	c, ok := httpClient.(*http.Client)
	if ok {
		c.Timeout = DefaultTimeout
	}
	client := &Client{
		entryPoint:       entryPoint,
		client:           httpClient,
		timeout:          DefaultTimeout,
		retryCount:       1,
		retryWaitTime:    DefaultWaitTime,
		retryMaxWaitTime: DefaultMaxWaitTime,
		maxResponseSize:  0,
		parseResponse:    DefaultParseResponse,
		parseError:       DefaultParseError,
		middleware:       NewClientMiddleware(),
		redirectConfig: RedirectConfig{
			policy:       FollowRedirects,
			maxRedirects: 0,
		},
		defaultHeaders: map[string]string{
			UserAgent: DefaultUserAgent,
		},
		mutex: &sync.RWMutex{},
	}
	for _, opt := range opts {
		opt(client)
	}

	return client
}

// SetTimeout sets the timeout for HTTP requests.
func (c *Client) SetTimeout(timeout time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.timeout = timeout
}

// SetEntryPoint sets the base URL for requests.
func (c *Client) SetEntryPoint(entryPoint string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.entryPoint = entryPoint
}

// SetRetryCount sets the number of retries on failure.
func (c *Client) SetRetryCount(count int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.retryCount = count
}

// SetRetryWaitTime sets the wait time between retries.
func (c *Client) SetRetryWaitTime(wait time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.retryWaitTime = wait
}

// SetRetryMaxWaitTime sets the maximum wait time between retries.
func (c *Client) SetRetryMaxWaitTime(wait time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.retryMaxWaitTime = wait
}

// SetParseResponse sets the function to parse successful responses.
func (c *Client) SetParseResponse(parseResponse ParseResponse) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.parseResponse = parseResponse
}

// SetParseError sets the function to parse error responses.
func (c *Client) SetParseError(parseError ParseResponse) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.parseError = parseError
}

// SetMaxResponseSize sets the maximum response body size in bytes.
// If the response body exceeds this size, it will be truncated.
func (c *Client) SetMaxResponseSize(size int64) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.maxResponseSize = size
}

// SetHeader sets a default header to be sent with all requests.
func (c *Client) SetHeader(header, value string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.defaultHeaders == nil {
		c.defaultHeaders = make(map[string]string)
	}
	c.defaultHeaders[header] = value
}

// SetHeaders sets multiple default headers to be sent with all requests.
func (c *Client) SetHeaders(headers map[string]string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.defaultHeaders == nil {
		c.defaultHeaders = make(map[string]string)
	}
	for k, v := range headers {
		c.defaultHeaders[k] = v
	}
}

// SetContentType sets the default Content-Type header for all requests.
func (c *Client) SetContentType(contentType string) {
	c.SetHeader(ContentType, contentType)
}

// UseMiddleware adds middleware to be executed before each request.
func (c *Client) UseMiddleware(middleware ...Middleware) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.middleware.Use(middleware...)
}

// SetRedirectPolicy sets the redirect policy (FollowRedirects or NoRedirect).
func (c *Client) SetRedirectPolicy(policy RedirectPolicy) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.redirectConfig.policy = policy
}

// SetMaxRedirects sets the maximum number of redirects to follow.
func (c *Client) SetMaxRedirects(maximum int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.redirectConfig.maxRedirects = maximum
}

// SetTransport sets the transport layer for the client.
func (c *Client) SetTransport(transport *http.Transport) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	httpClient, ok := c.client.(*http.Client)
	if ok {
		httpClient.Transport = transport
	}
}

// Execute sends the HTTP request and returns the response.
// It uses context.Background() for the request context.
func (c *Client) Execute(request *Request) (*Response, error) {
	return c.ExecuteWithContext(context.Background(), request)
}

// ExecuteWithContext sends the HTTP request with the given context and returns the response.
func (c *Client) ExecuteWithContext(ctx context.Context, request *Request) (*Response, error) {
	return c.middleware.Execute(request, func(req *Request) (*Response, error) {
		return c.doExecuteWithContext(ctx, req)
	})
}

func (c *Client) doExecuteWithContext(ctx context.Context, request *Request) (*Response, error) {
	var (
		req  *http.Request
		resp *http.Response
		err  error
	)

	c.mutex.RLock()
	entryPoint := c.entryPoint
	retryCount := c.retryCount
	retryWaitTime := min(c.retryWaitTime, c.retryMaxWaitTime)
	retryMaxWaitTime := c.retryMaxWaitTime
	maxResponseSize := c.maxResponseSize
	clientTimeout := c.timeout
	redirectConfig := c.redirectConfig
	c.mutex.RUnlock()

	timeout := clientTimeout
	if request.timeout > 0 {
		timeout = request.timeout
	}
	ctx, cancel := c.contextWithTimeout(ctx, timeout)
	if cancel != nil {
		defer cancel()
	}

	client := wrapWithRedirectPolicy(c.client, redirectConfig)

	for i := 0; i <= retryCount; i++ {
		req, err = request.computeWithContext(ctx, entryPoint, c.defaultHeaders)
		if err != nil {
			return nil, err
		}

		resp, err = client.Do(req)
		if err == nil {
			break
		}
		if !isRetriableError(err) {
			return nil, fmt.Errorf("%w: %w", ErrHTTPRequest, err)
		}
		time.Sleep(retryWaitTime)
		retryWaitTime = min(2*retryWaitTime, retryMaxWaitTime)
	}
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrHTTPRequest, err)
	}

	response := NewResponse(resp)
	if resp.Body == nil {
		return response, nil
	}
	defer resp.Body.Close()

	var reader io.Reader = resp.Body
	if maxResponseSize > 0 {
		reader = io.LimitReader(resp.Body, maxResponseSize)
	}

	response.bodyBytes, err = io.ReadAll(reader)
	if err != nil {
		return response, fmt.Errorf("%w: %w", ErrReadBody, err)
	}
	if response.IsError() && request.errorRespType != nil && c.parseError != nil {
		response.content, err = c.parseError(request, response)
		if err != nil {
			return response, err
		}
	}
	if !response.IsError() && request.respType != nil && c.parseResponse != nil {
		response.content, err = c.parseResponse(request, response)
		if err != nil {
			return response, err
		}
	}

	return response, nil
}

func (c *Client) contextWithTimeout(ctx context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	var cancel context.CancelFunc
	if timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, timeout)
	}
	return ctx, cancel
}

func isRetriableError(err error) bool {
	var (
		netErr net.Error
		urlErr *url.Error
	)
	if errors.Is(err, context.Canceled) {
		return false
	}

	if errors.Is(err, context.DeadlineExceeded) {
		return true
	}

	if errors.As(err, &netErr) && netErr.Timeout() {
		return true
	}

	if errors.As(err, &urlErr) {
		if urlErr.Op == "parse" {
			return false
		}
		err := urlErr.Err
		if err != nil && strings.Contains(err.Error(), "unsupported protocol scheme") {
			return false
		}
	}

	return true
}
