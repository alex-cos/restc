package restc

import (
	"context"
	"net"
	"net/http"
)

var zeroDialer net.Dialer

func NewIpv4Transport() *http.Transport {
	transport := http.DefaultTransport.(*http.Transport).Clone() // nolint: forcetypeassert
	transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		return zeroDialer.DialContext(ctx, "tcp4", addr)
	}

	return transport
}

func NewIpv6Transport() *http.Transport {
	transport := http.DefaultTransport.(*http.Transport).Clone() // nolint: forcetypeassert
	transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		return zeroDialer.DialContext(ctx, "tcp6", addr)
	}

	return transport
}
