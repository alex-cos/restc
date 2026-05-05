package restc

import (
	"context"
	"net"
	"net/http"
)

var zeroDialer net.Dialer

// NewIpv4Transport creates an http.Transport that uses only IPv4 connections.
func NewIpv4Transport() *http.Transport {
	transport := http.DefaultTransport.(*http.Transport).Clone() // nolint: forcetypeassert
	transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		return zeroDialer.DialContext(ctx, "tcp4", addr)
	}

	return transport
}

// NewIpv6Transport creates an http.Transport that uses only IPv6 connections.
func NewIpv6Transport() *http.Transport {
	transport := http.DefaultTransport.(*http.Transport).Clone() // nolint: forcetypeassert
	transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		return zeroDialer.DialContext(ctx, "tcp6", addr)
	}

	return transport
}
