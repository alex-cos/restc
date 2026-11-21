package restc

import (
	"context"
	"net"
	"net/http"
)

var zeroDialer net.Dialer

// cloneDefaultTransport returns a clone of http.DefaultTransport.
func cloneDefaultTransport() *http.Transport {
	if t, ok := http.DefaultTransport.(*http.Transport); ok {
		return t.Clone()
	}

	return &http.Transport{}
}

// NewIpv4Transport creates an http.Transport that uses only IPv4 connections.
func NewIpv4Transport(transports ...*http.Transport) *http.Transport {
	var transport *http.Transport

	if len(transports) > 0 {
		transport = transports[0]
	}
	if transport == nil {
		transport = cloneDefaultTransport()
	}
	transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		return zeroDialer.DialContext(ctx, "tcp4", addr)
	}

	return transport
}

// NewIpv6Transport creates an http.Transport that uses only IPv6 connections.
func NewIpv6Transport(transports ...*http.Transport) *http.Transport {
	var transport *http.Transport

	if len(transports) > 0 {
		transport = transports[0]
	}
	if transport == nil {
		transport = cloneDefaultTransport()
	}
	transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		return zeroDialer.DialContext(ctx, "tcp6", addr)
	}

	return transport
}
