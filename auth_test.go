package restc_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/alex-cos/restc"
	"github.com/stretchr/testify/assert"
)

func TestGetWithAuth(t *testing.T) {
	t.Parallel()

	client := restc.NewWithClient("",
		NewMockClient(http.StatusOK, GetResponse))
	client.SetEntryPoint("https://api.test.com")

	req := restc.Get("users").
		SetHeader("Accept", "application/json").
		SetAuthToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOjEsImlhdCI6MTY3Mjc2NjAyOCwiZXhwIjoxNjc0NDk0MDI4fQ.kCak9sLJr74frSRVQp0_27BY4iBCgQSmoT3vQVWKzJg").
		SetResponseType(&[]DummyObject{}).
		SetErrorRespType(&ReturnedError{})

	if !testing.Short() {
		fmt.Printf("req = %+v\n", req)
	}

	resp, err := client.Execute(req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.False(t, resp.IsError())
	assert.Equal(t, "HTTP/2.0", resp.Proto())
	assert.Equal(t, 200, resp.StatusCode())
	assert.Equal(t, "200 OK", resp.Status())
	assert.Equal(t, "application/json", resp.ContentType())
	assert.NotZero(t, resp.ReceivedAt())
	assert.NotEmpty(t, resp.Bytes())
	assert.NotEmpty(t, resp.Content())

	if !testing.Short() {
		fmt.Printf("resp = %+v\n", resp)
	}
}
