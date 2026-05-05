package restc

import (
	"net/http"
	"time"
)

type Option func(*Client)

func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.timeout = timeout
		if c.client == nil {
			return
		}
		httpClient, ok := c.client.(*http.Client)
		if ok {
			httpClient.Timeout = timeout
		}
	}
}

func WithRetryCount(count int) Option {
	return func(c *Client) {
		c.retryCount = count
	}
}

func WithRetryWaitTime(wait time.Duration) Option {
	return func(c *Client) {
		c.retryWaitTime = wait
	}
}

func WithRetryMaxWaitTime(wait time.Duration) Option {
	return func(c *Client) {
		c.retryMaxWaitTime = wait
	}
}

func WithParseResponse(parseResponse ParseResponse) Option {
	return func(c *Client) {
		c.parseResponse = parseResponse
	}
}

func WithParseError(parseError ParseResponse) Option {
	return func(c *Client) {
		c.parseError = parseError
	}
}

func WithHeader(header, value string) Option {
	return func(c *Client) {
		if c.defaultHeaders == nil {
			c.defaultHeaders = make(map[string]string)
		}
		c.defaultHeaders[header] = value
	}
}

func WithHeaders(headers map[string]string) Option {
	return func(c *Client) {
		if c.defaultHeaders == nil {
			c.defaultHeaders = make(map[string]string)
		}
		for k, v := range headers {
			c.defaultHeaders[k] = v
		}
	}
}

func WithRedirectPolicy(policy RedirectPolicy) Option {
	return func(c *Client) {
		c.redirectConfig.policy = policy
	}
}

func WithMaxRedirects(maximum int) Option {
	return func(c *Client) {
		c.redirectConfig.maxRedirects = maximum
	}
}

func WithMaxResponseSize(size int64) Option {
	return func(c *Client) {
		c.maxResponseSize = size
	}
}

func WithDisableIPv6() Option {
	return func(c *Client) {
		if c.client == nil {
			return
		}
		httpClient, ok := c.client.(*http.Client)
		if ok {
			httpClient.Transport = NewIpv4Transport()
		}
	}
}

func WithOnlyIPv6() Option {
	return func(c *Client) {
		if c.client == nil {
			return
		}
		httpClient, ok := c.client.(*http.Client)
		if ok {
			httpClient.Transport = NewIpv6Transport()
		}
	}
}
