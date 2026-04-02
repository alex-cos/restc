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
	DefaultTimeout     = 12 * time.Second
	DefaultWaitTime    = time.Duration(100) * time.Millisecond
	DefaultMaxWaitTime = time.Duration(2000) * time.Millisecond
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type ParseResponse func(request *Request, response *Response) (any, error)

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

func New(entryPoint string) *Client {
	return NewWithClient(entryPoint, http.DefaultClient)
}

func NewWithClient(entryPoint string, client HTTPClient) *Client {
	return NewWithClientTimeout(entryPoint, client, DefaultTimeout)
}

func NewWithTimeout(entryPoint string, timeout time.Duration) *Client {
	return NewWithClientTimeout(entryPoint, http.DefaultClient, timeout)
}

func NewWithClientTimeout(entryPoint string, httpClient HTTPClient, timeout time.Duration) *Client {
	return &Client{
		entryPoint:       entryPoint,
		client:           httpClient,
		timeout:          timeout,
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
		defaultHeaders: map[string]string{},
		mutex:          &sync.RWMutex{},
	}
}

func (c *Client) SetTimeout(timeout time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.timeout = timeout
}

func (c *Client) SetEntryPoint(entryPoint string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.entryPoint = entryPoint
}

func (c *Client) SetRetryCount(count int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.retryCount = count
}

func (c *Client) SetRetryWaitTime(wait time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.retryWaitTime = wait
}

func (c *Client) SetRetryMaxWaitTime(wait time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.retryMaxWaitTime = wait
}

func (c *Client) SetParseResponse(parseResponse ParseResponse) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.parseResponse = parseResponse
}

func (c *Client) SetParseError(parseError ParseResponse) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.parseError = parseError
}

func (c *Client) SetMaxResponseSize(size int64) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.maxResponseSize = size
}

func (c *Client) SetHeader(header, value string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.defaultHeaders == nil {
		c.defaultHeaders = make(map[string]string)
	}
	c.defaultHeaders[header] = value
}

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

func (c *Client) UseMiddleware(middleware ...Middleware) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.middleware.Use(middleware...)
}

func (c *Client) SetRedirectPolicy(policy RedirectPolicy) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.redirectConfig.policy = policy
}

func (c *Client) SetMaxRedirects(maximum int) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.redirectConfig.maxRedirects = maximum
}

func (c *Client) Execute(request *Request) (*Response, error) {
	return c.ExecuteWithContext(context.Background(), request)
}

func (c *Client) ExecuteWithContext(ctx context.Context, request *Request) (*Response, error) {
	return c.middleware.Execute(request, func(req *Request) (*Response, error) {
		return c.doExecuteWithContext(ctx, req)
	})
}

func (c *Client) doExecuteWithContext(ctx context.Context, request *Request) (*Response, error) {
	var (
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

	req, err := request.computeWithContext(ctx, entryPoint, c.defaultHeaders)
	if err != nil {
		return nil, err
	}

	for i := 0; i <= retryCount; i++ {
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
