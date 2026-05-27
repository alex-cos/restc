package restc

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"net/url"
	"os"
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

// WithContentType sets the default Content-Type header for all requests.
func WithContentType(contentType string) Option {
	return func(c *Client) {
		if c.defaultHeaders == nil {
			c.defaultHeaders = make(map[string]string)
		}
		c.defaultHeaders[ContentType] = contentType
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

// WithOnlyIPv4 configures the client to use only IPv4.
func WithOnlyIPv4() Option {
	return func(c *Client) {
		transport, httpClient, ok := cloneClientTransport(c)
		if !ok {
			return
		}
		httpClient.Transport = NewIpv4Transport(transport)
	}
}

// WithOnlyIPv6 configures the client to use only IPv6.
func WithOnlyIPv6() Option {
	return func(c *Client) {
		transport, httpClient, ok := cloneClientTransport(c)
		if !ok {
			return
		}
		httpClient.Transport = NewIpv6Transport(transport)
	}
}

// WithTLSConfig accepts a complete TLS configuration
// with preserving the actual transport layer.
func WithTLSConfig(config *tls.Config) Option {
	return func(c *Client) {
		transport, httpClient, ok := cloneClientTransport(c)
		if !ok {
			return
		}
		transport.TLSClientConfig = config
		httpClient.Transport = transport
	}
}

// WithMTLS configures the mutual TLS authentication
// with preserving the actual transport layer.
func WithMTLS(caCertFile, certFile, keyFile string) Option {
	return func(c *Client) {
		transport, httpClient, ok := cloneClientTransport(c)
		if !ok {
			return
		}
		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			return
		}
		caCertPEM, err := os.ReadFile(caCertFile)
		if err != nil {
			return
		}
		caPool := x509.NewCertPool()
		if !caPool.AppendCertsFromPEM(caCertPEM) {
			return
		}
		transport.TLSClientConfig = &tls.Config{
			RootCAs:      caPool,
			Certificates: []tls.Certificate{cert},
			MinVersion:   tls.VersionTLS12,
		}
		httpClient.Transport = transport
	}
}

// WithProxy configures the client to use the specified HTTP proxy.
// If user and password are provided, basic auth is added to the proxy URL.
func WithProxy(proxyURL, user, password string) Option {
	return func(c *Client) {
		transport, httpClient, ok := cloneClientTransport(c)
		if !ok {
			return
		}

		parsed, err := url.Parse(proxyURL)
		if err != nil {
			return
		}

		if user != "" && password != "" {
			parsed.User = url.UserPassword(user, password)
		} else if user != "" {
			parsed.User = url.User(user)
		}
		transport.Proxy = http.ProxyURL(parsed)
		httpClient.Transport = transport
	}
}

// WithTransportPool configures the HTTP transport connection pool settings.
// maxIdleConnsPerHost: max idle connections per host (default: 2)
// maxIdleConns: total max idle connections (default: 100)
// maxConnsPerHost: max connections per host, 0 = unlimited
// idleConnTimeout: how long idle connections stay in the pool.
func WithTransportPool(
	maxIdleConnsPerHost int,
	maxIdleConns int,
	maxConnsPerHost int,
	idleConnTimeout time.Duration,
) Option {
	return func(c *Client) {
		transport, httpClient, ok := cloneClientTransport(c)
		if !ok {
			return
		}

		transport.MaxIdleConnsPerHost = maxIdleConnsPerHost
		transport.MaxIdleConns = maxIdleConns
		transport.MaxConnsPerHost = maxConnsPerHost
		transport.IdleConnTimeout = idleConnTimeout

		httpClient.Transport = transport
	}
}

func cloneClientTransport(c *Client) (*http.Transport, *http.Client, bool) {
	if c.client == nil {
		return nil, nil, false
	}
	httpClient, ok := c.client.(*http.Client)
	if !ok {
		return nil, nil, false
	}
	// nolint: forcetypeassert
	transport := http.DefaultTransport.(*http.Transport).Clone()
	if t, ok := httpClient.Transport.(*http.Transport); ok {
		transport = t
	}
	return transport, httpClient, true
}
