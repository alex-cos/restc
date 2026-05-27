package restc_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/alex-cos/restc"
	"github.com/stretchr/testify/assert"
)

const GetResponse = `[
		{
			"id": 1,
			"firstname": "Emma",
			"lastname": "Bowen"
		},
		{
			"id": 2,
			"firstname": "Kevin",
			"lastname": "Banks"
		},
		{
			"id": 3,
			"firstname": "Paul",
			"lastname": "Wang"
		},
    {
			"id": 4,
			"firstname": "Catherine",
			"lastname": "Nichols"
		}
	]`

func TestGetSucess(t *testing.T) {
	t.Parallel()

	client := restc.NewWithClient("https://api.test.com",
		NewMockClient(http.StatusOK, GetResponse),
		restc.WithTimeout(5*time.Second),
		restc.WithRetryCount(3),
		restc.WithRetryWaitTime(200*time.Millisecond),
		restc.WithRetryMaxWaitTime(2*time.Second),
		restc.WithMaxResponseSize(10*1024*1024),
		restc.WithRedirectPolicy(restc.NoRedirect),
		restc.WithContentType(restc.TypeApplicationJSON),
		restc.WithHeader("User-Agent", "MyAgent/1.0"),
	)

	req := restc.NewRequest("", "users").
		SetMethod(restc.MethodGet).
		SetHeader("Accept", restc.TypeApplicationJSON)

	if !testing.Short() {
		fmt.Printf("req = %+v\n", req)
	}

	resp, err := client.Execute(req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "HTTP/2.0", resp.Proto())
	assert.Equal(t, 200, resp.StatusCode())
	assert.Equal(t, "200 OK", resp.Status())
	assert.Contains(t, resp.ContentType(), restc.TypeApplicationJSON)
	assert.NotZero(t, resp.ReceivedAt())
	assert.NotEmpty(t, resp.Bytes())
	assert.Nil(t, resp.Content())

	if !testing.Short() {
		fmt.Printf("resp = %+v\n", resp)
	}
}

func TestGetWithURL(t *testing.T) {
	t.Parallel()

	client := restc.NewWithClient("https://api.test.com",
		NewMockClient(http.StatusOK, GetResponse),
		restc.WithTimeout(5*time.Second),
		restc.WithRetryCount(3),
		restc.WithRetryWaitTime(200*time.Millisecond),
		restc.WithRetryMaxWaitTime(2*time.Second),
		restc.WithMaxResponseSize(10*1024*1024),
		restc.WithRedirectPolicy(restc.NoRedirect),
		restc.WithContentType(restc.TypeApplicationJSON),
		restc.WithHeaders(map[string]string{"User-Agent": "MyAgent/1.0"}),
	)

	req := restc.Get("users").
		SetMethod(restc.MethodGet).
		SetHeader("Accept", restc.TypeApplicationJSON)
	req.SetURL("https://api.test.com/v2/users")

	if !testing.Short() {
		fmt.Printf("req = %+v\n", req)
	}

	resp, err := client.Execute(req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "HTTP/2.0", resp.Proto())
	assert.Equal(t, 200, resp.StatusCode())
	assert.Equal(t, "200 OK", resp.Status())
	assert.Contains(t, resp.ContentType(), restc.TypeApplicationJSON)
	assert.NotZero(t, resp.ReceivedAt())
	assert.NotEmpty(t, resp.Bytes())
	assert.Nil(t, resp.Content())

	if !testing.Short() {
		fmt.Printf("resp = %+v\n", resp)
	}
}

func TestGetWithCookies(t *testing.T) {
	t.Parallel()

	client := restc.NewWithClient("https://api.test.com",
		NewMockClient(http.StatusOK, GetResponse))

	req := restc.NewRequest("", "users").
		SetMethod(restc.MethodGet).
		SetCookie(&http.Cookie{
			Name:     "test",
			Value:    "ok",
			Secure:   true,
			HttpOnly: true,
		}).
		SetHeader("Accept", restc.TypeApplicationJSON)

	if !testing.Short() {
		fmt.Printf("req = %+v\n", req)
	}

	resp, err := client.Execute(req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "HTTP/2.0", resp.Proto())
	assert.Equal(t, 200, resp.StatusCode())
	assert.Equal(t, "200 OK", resp.Status())
	assert.Contains(t, resp.ContentType(), restc.TypeApplicationJSON)
	assert.NotZero(t, resp.ReceivedAt())
	assert.NotEmpty(t, resp.Bytes())
	assert.Nil(t, resp.Content())

	if !testing.Short() {
		fmt.Printf("resp = %+v\n", resp)
	}
}

func TestHeadSucess(t *testing.T) {
	t.Parallel()

	client := restc.NewWithClient("https://api.test.com",
		NewMockClient(http.StatusOK, GetResponse))

	req := restc.Head("users").
		SetMethod(restc.MethodHead).
		SetHeader("Accept", restc.TypeApplicationJSON)

	if !testing.Short() {
		fmt.Printf("req = %+v\n", req)
	}

	resp, err := client.Execute(req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "HTTP/2.0", resp.Proto())
	assert.Equal(t, 200, resp.StatusCode())
	assert.Equal(t, "200 OK", resp.Status())
	assert.Nil(t, resp.Content())

	if !testing.Short() {
		fmt.Printf("resp = %+v\n", resp)
	}
}

func TestOptionsSucess(t *testing.T) {
	t.Parallel()

	client := restc.NewWithClient("https://api.test.com",
		NewMockClient(http.StatusOK, GetResponse))

	req := restc.Options("users").
		SetHeader("Accept", restc.TypeApplicationJSON)

	if !testing.Short() {
		fmt.Printf("req = %+v\n", req)
	}

	resp, err := client.Execute(req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "HTTP/2.0", resp.Proto())
	assert.Equal(t, 200, resp.StatusCode())
	assert.Equal(t, "200 OK", resp.Status())
	assert.Nil(t, resp.Content())

	if !testing.Short() {
		fmt.Printf("resp = %+v\n", resp)
	}
}

func TestTraceSucess(t *testing.T) {
	t.Parallel()

	client := restc.NewWithClient("https://api.test.com",
		NewMockClient(http.StatusOK, GetResponse))

	req := restc.Trace("users").
		SetHeader("Accept", restc.TypeApplicationJSON)

	if !testing.Short() {
		fmt.Printf("req = %+v\n", req)
	}

	resp, err := client.Execute(req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "HTTP/2.0", resp.Proto())
	assert.Equal(t, 200, resp.StatusCode())
	assert.Equal(t, "200 OK", resp.Status())
	assert.Nil(t, resp.Content())

	if !testing.Short() {
		fmt.Printf("resp = %+v\n", resp)
	}
}

func TestGetWithParamsSucess(t *testing.T) {
	t.Parallel()

	client := restc.NewWithClient("https://api.test.com",
		NewMockClient(http.StatusOK, `[
		{
			"id": 2,
			"firstname": "Kevin",
			"lastname": "Banks",
		},
    {
			"id": 4,
			"firstname": "Catherine",
			"lastname": "Nichols",
		},
	]`))

	req := restc.Get("users").
		AddQueryParam("id", "2").
		AddQueryParam("id", "4").
		SetHeader("Accept", restc.TypeApplicationJSON)

	if !testing.Short() {
		fmt.Printf("req = %+v\n", req)
	}

	resp, err := client.Execute(req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "HTTP/2.0", resp.Proto())
	assert.Equal(t, 200, resp.StatusCode())
	assert.Equal(t, "200 OK", resp.Status())
	assert.Contains(t, resp.ContentType(), restc.TypeApplicationJSON)
	assert.NotZero(t, resp.ReceivedAt())
	assert.NotEmpty(t, resp.Bytes())
	assert.Nil(t, resp.Content())

	if !testing.Short() {
		fmt.Printf("resp = %+v\n", resp)
	}
}

func TestPostSucess(t *testing.T) {
	t.Parallel()

	client := restc.NewWithClient("https://api.test.com",
		NewMockClient(http.StatusOK, `{
			"id": 5,
			"firstname": "John",
			"lastname": "Doe"
		}`))
	client.SetTimeout(5 * time.Second)

	req := restc.Post("users").
		SetHeader("Accept", restc.TypeApplicationJSON).
		SetContentType(restc.TypeApplicationJSON).
		SetBody(`{
			"firstname": "John",
			"lastname": "Doe"
		}`)

	if !testing.Short() {
		fmt.Printf("req = %+v\n", req)
	}

	resp, err := client.Execute(req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "HTTP/2.0", resp.Proto())
	assert.Equal(t, 200, resp.StatusCode())
	assert.Equal(t, "200 OK", resp.Status())
	assert.Contains(t, resp.ContentType(), restc.TypeApplicationJSON)
	assert.NotZero(t, resp.ReceivedAt())
	assert.NotEmpty(t, resp.Bytes())
	assert.Nil(t, resp.Content())

	if !testing.Short() {
		fmt.Printf("resp = %+v\n", resp)
	}

	var data struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
	err = json.Unmarshal(resp.Bytes(), &data)
	assert.NoError(t, err)

	if !testing.Short() {
		fmt.Printf("data = %+v\n", data)
	}
}

func TestUpdateSucess(t *testing.T) {
	t.Parallel()

	client := restc.NewWithClient("https://api.test.com",
		NewMockClient(http.StatusOK, `{
			"id": 3,
			"firstname": "Paul",
			"lastname": "Klein"
		}`),
		restc.WithTimeout(5*time.Second))

	req := restc.Put("users/3").
		SetHeader("Accept", restc.TypeApplicationJSON).
		SetBody(`{
			"firstname": "Paul",
			"lastname": "Klein"
		}`)

	if !testing.Short() {
		fmt.Printf("req = %+v\n", req)
	}

	resp, err := client.Execute(req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "HTTP/2.0", resp.Proto())
	assert.Equal(t, 200, resp.StatusCode())
	assert.Equal(t, "200 OK", resp.Status())
	assert.Contains(t, resp.ContentType(), restc.TypeApplicationJSON)
	assert.NotZero(t, resp.ReceivedAt())
	assert.NotEmpty(t, resp.Bytes())
	assert.Nil(t, resp.Content())

	if !testing.Short() {
		fmt.Printf("resp = %+v\n", resp)
	}
}

func TestPatchSucess(t *testing.T) {
	t.Parallel()

	client := restc.NewWithClient("https://api.test.com",
		NewMockClient(http.StatusOK, `{
			"id": 3,
			"firstname": "Paul",
			"lastname": "Klein"
		}`),
		restc.WithTimeout(5*time.Second))

	req := restc.Patch("users/3").
		SetHeader("Accept", restc.TypeApplicationJSON).
		SetBody(`{"lastname": "Klein"}`)

	if !testing.Short() {
		fmt.Printf("req = %+v\n", req)
	}

	resp, err := client.Execute(req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "HTTP/2.0", resp.Proto())
	assert.Equal(t, 200, resp.StatusCode())
	assert.Equal(t, "200 OK", resp.Status())
	assert.Contains(t, resp.ContentType(), restc.TypeApplicationJSON)
	assert.NotZero(t, resp.ReceivedAt())
	assert.NotEmpty(t, resp.Bytes())
	assert.Nil(t, resp.Content())

	if !testing.Short() {
		fmt.Printf("resp = %+v\n", resp)
	}
}

func TestDeleteSucess(t *testing.T) {
	t.Parallel()

	client := restc.NewWithClient("https://api.test.com",
		NewMockClient(http.StatusOK, `{
			"message": "User with id = 3 has been successfully deleted."
		}`),
		restc.WithTimeout(5*time.Second))

	req := restc.Delete("users/3").
		SetHeader("Accept", restc.TypeApplicationJSON)

	if !testing.Short() {
		fmt.Printf("req = %+v\n", req)
	}

	resp, err := client.Execute(req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "HTTP/2.0", resp.Proto())
	assert.Equal(t, 200, resp.StatusCode())
	assert.Equal(t, "200 OK", resp.Status())
	assert.Contains(t, resp.ContentType(), restc.TypeApplicationJSON)
	assert.NotZero(t, resp.ReceivedAt())
	assert.NotEmpty(t, resp.Bytes())
	assert.Nil(t, resp.Content())

	if !testing.Short() {
		fmt.Printf("resp = %+v\n", resp)
	}
}
