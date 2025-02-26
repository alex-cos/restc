package restc_test

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/alex-cos/restc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestWrongMethod(t *testing.T) {
	t.Parallel()

	client := restc.NewWithClient("https://api.test.com",
		NewMockClient(http.StatusOK, ""))

	req := restc.NewRequest("ABCD", "users").
		SetHeader("Accept", "application/json")

	if !testing.Short() {
		fmt.Printf("req = %+v\n", req)
	}

	resp, err := client.Execute(req)
	assert.ErrorContains(t, err, "invalid provided HTTP method")
	assert.Nil(t, resp)

	if !testing.Short() {
		fmt.Printf("resp = %+v\n", resp)
	}
}

func TestGetWithWrongType(t *testing.T) {
	t.Parallel()

	client := restc.NewWithClient("",
		NewMockClient(http.StatusOK, GetResponse))
	client.SetEntryPoint("https://api.test.com")

	req := restc.Get("users").
		SetHeader("Accept", "application/json").
		SetResponseType(&DummyObject{}).
		SetErrorRespType(&ReturnedError{})

	if !testing.Short() {
		fmt.Printf("req = %+v\n", req)
	}

	resp, err := client.Execute(req)
	assert.ErrorContains(t, err, "failed to parse JSON response")
	assert.NotNil(t, resp)
	assert.False(t, resp.IsError())
	assert.Equal(t, "HTTP/2.0", resp.Proto())
	assert.Equal(t, 200, resp.StatusCode())
	assert.Equal(t, "200 OK", resp.Status())
	assert.Equal(t, "application/json", resp.ContentType())
	assert.NotZero(t, resp.ReceivedAt())
	assert.NotEmpty(t, resp.Bytes())

	if !testing.Short() {
		fmt.Printf("resp = %+v\n", resp)
	}
}

func TestGetJSONErrorWithType(t *testing.T) {
	t.Parallel()

	client := restc.NewWithClient("https://api.test.com",
		NewMockClient(http.StatusInternalServerError, `{
		"status": 500,
		"error": "Internal Server Error",
		"path": "users"
	 }`))

	req := restc.Get("users").
		SetHeader("Accept", "application/json").
		SetResponseType(&[]DummyObject{}).
		SetErrorRespType(&ReturnedError{})

	if !testing.Short() {
		fmt.Printf("req = %+v\n", req)
	}

	resp, err := client.Execute(req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.IsError())
	assert.Equal(t, "HTTP/2.0", resp.Proto())
	assert.Equal(t, 500, resp.StatusCode())
	assert.Equal(t, "500 Internal Server Error", resp.Status())
	assert.Equal(t, "application/json", resp.ContentType())
	assert.NotZero(t, resp.ReceivedAt())
	assert.NotEmpty(t, resp.Bytes())
	assert.NotNil(t, resp.Content())
	assert.Equal(t, &ReturnedError{
		Status: 500,
		Error:  "Internal Server Error",
		Path:   "users",
	}, resp.Content())

	if !testing.Short() {
		fmt.Printf("resp = %+v\n", resp)
	}
}

func TestGetUnexpectedType(t *testing.T) {
	t.Parallel()

	body := io.NopCloser(bytes.NewReader([]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF}))

	httpclient := NewMockClient(http.StatusOK, "").(*Mock)
	httpclient.On("Do", mock.Anything).Unset()
	httpclient.On("Do", mock.Anything).Return(&http.Response{
		Status:           "200 OK",
		StatusCode:       http.StatusOK,
		Proto:            "HTTP/2.0",
		ProtoMajor:       2,
		ProtoMinor:       0,
		Header:           nil,
		Body:             body,
		ContentLength:    int64(5),
		TransferEncoding: []string{},
		Close:            false,
		Uncompressed:     false,
		Trailer:          http.Header{},
		Request:          nil,
		TLS:              nil,
	}, nil)

	client := restc.NewWithClient("https://api.test.com", httpclient)

	req := restc.Get("users").
		SetHeader("Accept", "application/json").
		SetResponseType(&[]DummyObject{}).
		SetErrorRespType(&ReturnedError{})

	if !testing.Short() {
		fmt.Printf("req = %+v\n", req)
	}

	resp, err := client.Execute(req)
	assert.ErrorContains(t, err, "unexpected response type")
	assert.NotNil(t, resp)
	assert.False(t, resp.IsError())
	assert.Equal(t, "HTTP/2.0", resp.Proto())
	assert.Equal(t, 200, resp.StatusCode())
	assert.Equal(t, "200 OK", resp.Status())
	assert.NotZero(t, resp.ReceivedAt())
	assert.NotEmpty(t, resp.Bytes())
	assert.Nil(t, resp.Content())

	if !testing.Short() {
		fmt.Printf("resp = %+v\n", resp)
	}
}

func TestGetHTMLErrorWithType(t *testing.T) {
	t.Parallel()

	testfile := filepath.Join("testdata", "test1.html")
	data, err := os.ReadFile(testfile)
	assert.NoError(t, err)

	httpclient := NewMockClient(http.StatusMethodNotAllowed, "")
	(httpclient.(*Mock)).SetContent(string(data))
	(httpclient.(*Mock)).SetContentType(restc.TypeTextHTML)

	client := restc.NewWithClient("https://api.test.com", httpclient)

	req := restc.Get("xxxxx").
		SetHeader("Accept", "application/json").
		SetResponseType(&[]DummyObject{}).
		SetErrorRespType(&ReturnedError{})

	if !testing.Short() {
		fmt.Printf("req = %+v\n", req)
	}

	resp, err := client.Execute(req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.IsError())
	assert.Equal(t, "HTTP/2.0", resp.Proto())
	assert.Equal(t, 405, resp.StatusCode())
	assert.Equal(t, "405 Method Not Allowed", resp.Status())
	assert.Equal(t, "text/html", resp.ContentType())
	assert.NotZero(t, resp.ReceivedAt())
	assert.NotEmpty(t, resp.Bytes())
	assert.NotNil(t, resp.Content())
	assert.Equal(t, "Unexpected Error: Something went wrong and cause an unexpected error.", resp.Content())

	if !testing.Short() {
		fmt.Printf("resp = %+v\n", resp)
	}
}

func TestGetWithWrongURL(t *testing.T) {
	t.Parallel()

	client := restc.New("https://wrong")
	client.SetTimeout(500 * time.Millisecond)
	client.SetRetryCount(2)
	client.SetRetryWaitTime(25 * time.Millisecond)
	client.SetRetryMaxWaitTime(time.Second)

	req := restc.Get("users").
		SetHeader("Accept", "application/json").
		SetResponseType(&[]DummyObject{}).
		SetErrorRespType(&ReturnedError{})

	if !testing.Short() {
		fmt.Printf("req = %+v\n", req)
	}

	resp, err := client.Execute(req)
	assert.Error(t, err)
	assert.Nil(t, resp)
}
