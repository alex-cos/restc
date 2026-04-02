package restc

import (
	"fmt"
	"net/http"
)

type RedirectPolicy int

const (
	FollowRedirects RedirectPolicy = iota
	NoRedirect
)

type RedirectConfig struct {
	policy       RedirectPolicy
	maxRedirects int
}

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
