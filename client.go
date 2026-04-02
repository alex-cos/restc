package restc

import (
	"context"
	"fmt"
	"io"
	"net/http"
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
	parseResponse    ParseResponse
	parseError       ParseResponse
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
		parseResponse:    DefaultParseResponse,
		parseError:       DefaultParseError,
		mutex:            &sync.RWMutex{},
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

func (c *Client) Execute(request *Request) (*Response, error) {
	return c.ExecuteWithContext(context.Background(), request)
}

func (c *Client) ExecuteWithContext(ctx context.Context, request *Request) (*Response, error) {
	var (
		resp *http.Response
		err  error
	)

	c.mutex.RLock()
	entryPoint := c.entryPoint
	retryCount := c.retryCount
	retryWaitTime := minDuration(c.retryWaitTime, c.retryMaxWaitTime)
	retryMaxWaitTime := c.retryMaxWaitTime
	c.mutex.RUnlock()

	ctx, cancel := c.context(ctx)
	if cancel != nil {
		defer cancel()
	}

	req, err := request.computeWithContext(ctx, entryPoint)
	if err != nil {
		return nil, err
	}

	for i := 0; i <= retryCount; i++ {
		resp, err = c.client.Do(req)
		if err == nil {
			break
		}
		time.Sleep(retryWaitTime)
		retryWaitTime = minDuration(2*retryWaitTime, retryMaxWaitTime)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to execute HTTP request: %w", err)
	}

	response := NewResponse(resp)
	if resp.Body == nil {
		return response, nil
	}
	defer resp.Body.Close()

	response.bodyBytes, err = io.ReadAll(resp.Body)
	if err != nil {
		return response, fmt.Errorf("failed to read body response: %w", err)
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

func (c *Client) context(ctx context.Context) (context.Context, context.CancelFunc) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	var cancel context.CancelFunc
	if c.timeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, c.timeout)
	}

	return ctx, cancel
}
