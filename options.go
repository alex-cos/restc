package restc

import (
	"net/http"
	"time"
)

// Option is a function that configures a Client.
// It is used with New or NewWithClient to configure the client.
type Option func(*Client)

// WithTimeout sets the timeout for HTTP requests.
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

// WithRetryCount sets the number of retries on failure.
func WithRetryCount(count int) Option {
	return func(c *Client) {
		c.retryCount = count
	}
}

// WithRetryWaitTime sets the initial wait time between retries.
func WithRetryWaitTime(wait time.Duration) Option {
	return func(c *Client) {
		c.retryWaitTime = wait
	}
}

// WithRetryMaxWaitTime sets the maximum wait time between retries.
func WithRetryMaxWaitTime(wait time.Duration) Option {
	return func(c *Client) {
		c.retryMaxWaitTime = wait
	}
}

// WithParseResponse sets the function to parse successful responses.
func WithParseResponse(parseResponse ParseResponse) Option {
	return func(c *Client) {
		c.parseResponse = parseResponse
	}
}

// WithParseError sets the function to parse error responses.
func WithParseError(parseError ParseResponse) Option {
	return func(c *Client) {
		c.parseError = parseError
	}
}

// WithHeader sets a default header to be sent with all requests.
func WithHeader(header, value string) Option {
	return func(c *Client) {
		if c.defaultHeaders == nil {
			c.defaultHeaders = make(map[string]string)
		}
		c.defaultHeaders[header] = value
	}
}

// WithHeaders sets multiple default headers to be sent with all requests.
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

// WithRedirectPolicy sets the redirect policy.
func WithRedirectPolicy(policy RedirectPolicy) Option {
	return func(c *Client) {
		c.redirectConfig.policy = policy
	}
}

// WithMaxRedirects sets the maximum number of redirects to follow.
func WithMaxRedirects(maximum int) Option {
	return func(c *Client) {
		c.redirectConfig.maxRedirects = maximum
	}
}

// WithMaxResponseSize sets the maximum response body size in bytes.
func WithMaxResponseSize(size int64) Option {
	return func(c *Client) {
		c.maxResponseSize = size
	}
}

// WithDisableIPv6 configures the client to use only IPv4.
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

// WithOnlyIPv6 configures the client to use only IPv6.
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
