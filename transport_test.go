package restc_test

import (
	"crypto/tls"
	"net/http"
	"testing"
	"time"

	"github.com/alex-cos/restc"
	"github.com/stretchr/testify/assert"
)

func TestNewIpv4Transport_Default(t *testing.T) {
	t.Parallel()

	transport := restc.NewIpv4Transport()
	assert.NotNil(t, transport)
	assert.NotNil(t, transport.DialContext)
}

func TestNewIpv4Transport_Custom(t *testing.T) {
	t.Parallel()

	customTransport := &http.Transport{
		TLSClientConfig: &tls.Config{
			ServerName: "test.example.com",
		},
	}

	result := restc.NewIpv4Transport(customTransport)

	assert.Same(t, customTransport, result)
	assert.NotNil(t, result.DialContext)
	assert.NotNil(t, result.TLSClientConfig)
	assert.Equal(t, "test.example.com", result.TLSClientConfig.ServerName)
}

func TestNewIpv6Transport_Default(t *testing.T) {
	t.Parallel()

	transport := restc.NewIpv6Transport()
	assert.NotNil(t, transport)
	assert.NotNil(t, transport.DialContext)
}

func TestNewIpv6Transport_Custom(t *testing.T) {
	t.Parallel()

	customTransport := &http.Transport{
		TLSClientConfig: &tls.Config{
			ServerName: "test.example.com",
		},
	}

	result := restc.NewIpv6Transport(customTransport)

	assert.Same(t, customTransport, result)
	assert.NotNil(t, result.DialContext)
	assert.NotNil(t, result.TLSClientConfig)
	assert.Equal(t, "test.example.com", result.TLSClientConfig.ServerName)
}

func TestWithOnlyIPv4(t *testing.T) {
	t.Parallel()

	httpClient := &http.Client{}
	restc.NewWithClient("https://api.test.com",
		httpClient,
		restc.WithOnlyIPv4(),
	)

	transport, ok := httpClient.Transport.(*http.Transport)
	assert.True(t, ok)
	assert.NotNil(t, transport)
	assert.NotNil(t, transport.DialContext)
}

func TestWithOnlyIPv4_PreservesExistingTransport(t *testing.T) {
	t.Parallel()

	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				ServerName: "test.example.com",
			},
		},
	}
	restc.NewWithClient("https://api.test.com",
		httpClient,
		restc.WithOnlyIPv4(),
	)

	transport, ok := httpClient.Transport.(*http.Transport)
	assert.True(t, ok)
	assert.NotNil(t, transport.DialContext)
	assert.NotNil(t, transport.TLSClientConfig)
	assert.Equal(t, "test.example.com", transport.TLSClientConfig.ServerName)
}

func TestWithOnlyIPv6(t *testing.T) {
	t.Parallel()

	httpClient := &http.Client{}
	restc.NewWithClient("https://api.test.com",
		httpClient,
		restc.WithOnlyIPv6(),
	)

	transport, ok := httpClient.Transport.(*http.Transport)
	assert.True(t, ok)
	assert.NotNil(t, transport)
	assert.NotNil(t, transport.DialContext)
}

func TestWithOnlyIPv6_PreservesExistingTransport(t *testing.T) {
	t.Parallel()

	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				ServerName: "test.example.com",
			},
		},
	}
	restc.NewWithClient("https://api.test.com",
		httpClient,
		restc.WithOnlyIPv6(),
	)

	transport, ok := httpClient.Transport.(*http.Transport)
	assert.True(t, ok)
	assert.NotNil(t, transport.DialContext)
	assert.NotNil(t, transport.TLSClientConfig)
	assert.Equal(t, "test.example.com", transport.TLSClientConfig.ServerName)
}

func TestWithProxy(t *testing.T) {
	t.Parallel()

	httpClient := &http.Client{}
	restc.NewWithClient("https://api.test.com",
		httpClient,
		restc.WithProxy("http://proxy.example.com:8080", "", ""),
	)

	transport, ok := httpClient.Transport.(*http.Transport)
	assert.True(t, ok)
	assert.NotNil(t, transport.Proxy)
}

func TestWithProxy_WithAuth(t *testing.T) {
	t.Parallel()

	httpClient := &http.Client{}
	restc.NewWithClient("https://api.test.com",
		httpClient,
		restc.WithProxy("http://proxy.example.com:8080", "admin", "secret"),
	)

	transport, ok := httpClient.Transport.(*http.Transport)
	assert.True(t, ok)
	assert.NotNil(t, transport.Proxy)

	// Verify the proxy URL contains the auth
	req, _ := http.NewRequest(http.MethodGet, "https://api.test.com", nil)
	proxyURL, err := transport.Proxy(req)
	assert.NoError(t, err)
	assert.Equal(t, "admin", proxyURL.User.Username())
	pass, _ := proxyURL.User.Password()
	assert.Equal(t, "secret", pass)
}

func TestWithProxy_InvalidURL(t *testing.T) {
	t.Parallel()

	httpClient := &http.Client{}
	restc.NewWithClient("https://api.test.com",
		httpClient,
		restc.WithProxy("://invalid-url", "", ""),
	)

	assert.Nil(t, httpClient.Transport)
}

func TestWithProxy_PreservesExistingTransport(t *testing.T) {
	t.Parallel()

	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				ServerName: "test.example.com",
			},
		},
	}
	restc.NewWithClient("https://api.test.com",
		httpClient,
		restc.WithProxy("http://proxy.example.com:8080", "", ""),
	)

	transport, ok := httpClient.Transport.(*http.Transport)
	assert.True(t, ok)
	assert.NotNil(t, transport.Proxy)
	assert.NotNil(t, transport.TLSClientConfig)
	assert.Equal(t, "test.example.com", transport.TLSClientConfig.ServerName)
}

func TestWithTransportPool(t *testing.T) {
	t.Parallel()

	httpClient := &http.Client{}
	restc.NewWithClient("https://api.test.com",
		httpClient,
		restc.WithTransportPool(10, 200, 50, 60*time.Second),
	)

	transport, ok := httpClient.Transport.(*http.Transport)
	assert.True(t, ok)
	assert.Equal(t, 10, transport.MaxIdleConnsPerHost)
	assert.Equal(t, 200, transport.MaxIdleConns)
	assert.Equal(t, 50, transport.MaxConnsPerHost)
	assert.Equal(t, 60*time.Second, transport.IdleConnTimeout)
}

func TestWithTransportPool_PreservesExistingTransport(t *testing.T) {
	t.Parallel()

	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				ServerName: "test.example.com",
			},
		},
	}
	restc.NewWithClient("https://api.test.com",
		httpClient,
		restc.WithTransportPool(10, 200, 50, 60*time.Second),
	)

	transport, ok := httpClient.Transport.(*http.Transport)
	assert.True(t, ok)
	assert.Equal(t, 10, transport.MaxIdleConnsPerHost)
	assert.Equal(t, 200, transport.MaxIdleConns)
	assert.Equal(t, 50, transport.MaxConnsPerHost)
	assert.Equal(t, 60*time.Second, transport.IdleConnTimeout)
	assert.NotNil(t, transport.TLSClientConfig)
	assert.Equal(t, "test.example.com", transport.TLSClientConfig.ServerName)
}
