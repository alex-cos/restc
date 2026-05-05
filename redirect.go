package restc

import (
	"fmt"
	"net/http"
)

// RedirectPolicy defines the behavior for handling HTTP redirects.
type RedirectPolicy int

const (
	// FollowRedirects indicates that the client should follow redirects.
	FollowRedirects RedirectPolicy = iota
	// NoRedirect indicates that the client should not follow redirects.
	NoRedirect
)

// RedirectConfig holds configuration for redirect handling.
type RedirectConfig struct {
	policy       RedirectPolicy
	maxRedirects int
}

// checkRedirect determines whether to follow a redirect.
// It is called by the HTTP client for each redirect.
func (rc RedirectConfig) checkRedirect(_ *http.Request, via []*http.Request) error {
	switch rc.policy {
	case NoRedirect:
		return http.ErrUseLastResponse
	case FollowRedirects:
		if rc.maxRedirects > 0 && len(via) >= rc.maxRedirects {
			return fmt.Errorf("%w: %w",
				ErrMaxRedirects,
				fmt.Errorf("stopped after %d redirects", rc.maxRedirects),
			)
		}
		return nil
	default:
		return nil
	}
}

func wrapWithRedirectPolicy(client HTTPClient, config RedirectConfig) HTTPClient {
	httpClient, ok := client.(*http.Client)
	if !ok {
		return client
	}

	cloned := *httpClient
	cloned.CheckRedirect = config.checkRedirect
	return &cloned
}
