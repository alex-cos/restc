package restc_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alex-cos/restc"
	"github.com/stretchr/testify/assert"
)

func TestRedirectFollowsByDefault(t *testing.T) {
	t.Parallel()

	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`)) //nolint: errcheck
	}))
	defer target.Close()

	redirector := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, target.URL, http.StatusFound)
	}))
	defer redirector.Close()

	client := restc.NewWithClient(redirector.URL, &http.Client{})

	resp, err := client.Execute(restc.Get("/"))
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode())
}

func TestRedirectNoRedirect(t *testing.T) {
	t.Parallel()

	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`)) //nolint: errcheck
	}))
	defer target.Close()

	redirector := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, target.URL, http.StatusFound)
	}))
	defer redirector.Close()

	client := restc.NewWithClient(redirector.URL, &http.Client{})
	client.SetRedirectPolicy(restc.NoRedirect)

	resp, err := client.Execute(restc.Get("/"))
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusFound, resp.StatusCode())
}
