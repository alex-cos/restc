package restc_test

import (
	"net/http"
	"testing"

	"github.com/alex-cos/restc"
	"github.com/stretchr/testify/assert"
)

func TestRequestClone(t *testing.T) {
	t.Parallel()

	baseReq := restc.Get("users").
		SetHeader("Accept", restc.TypeApplicationJSON).
		SetQueryParam("limit", "10").
		SetAuthToken("token123")

	clonedReq := baseReq.Clone()

	clonedReq.SetQueryParam("limit", "20").
		SetHeader("X-Request-ID", "abc")

	assert.Contains(t, baseReq.String(), "limit=10")
	assert.Contains(t, clonedReq.String(), "limit=20")
}

func TestRequestCloneIndependence(t *testing.T) {
	t.Parallel()

	client := restc.NewWithClient("https://api.test.com",
		NewMockClient(http.StatusOK, GetResponse))

	baseReq := restc.Get("users").
		SetHeader("Accept", restc.TypeApplicationJSON)

	clonedReq := baseReq.Clone()
	clonedReq.SetHeader("Accept", "application/xml")

	resp, err := client.Execute(baseReq)
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	resp2, err := client.Execute(clonedReq)
	assert.NoError(t, err)
	assert.NotNil(t, resp2)
}
