package restc_test

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"testing"

	"github.com/alex-cos/restc"
	"github.com/stretchr/testify/assert"
)

type multipartMock struct {
	lastContentType string
	onDo            func(req *http.Request) (*http.Response, error)
}

func (m *multipartMock) Do(req *http.Request) (*http.Response, error) {
	m.lastContentType = req.Header.Get("Content-Type")
	if m.onDo != nil {
		return m.onDo(req)
	}
	return m.mockResponse()
}

func (m *multipartMock) mockResponse() (*http.Response, error) {
	body := io.NopCloser(strings.NewReader(`{"status": "ok"}`))
	return &http.Response{
		Status:     "200 OK",
		StatusCode: http.StatusOK,
		Proto:      "HTTP/2.0",
		Header: http.Header{
			"Content-Type": []string{"application/json"},
		},
		Body: body,
	}, nil
}

func TestMultipartFormDataOnly(t *testing.T) {
	t.Parallel()

	var capturedBody []byte

	mockClient := &multipartMock{}
	mockClient.onDo = func(req *http.Request) (*http.Response, error) {
		body, _ := io.ReadAll(req.Body)
		capturedBody = body
		return mockClient.mockResponse()
	}

	client := restc.NewWithClient("https://api.test.com", mockClient)

	req := restc.Post("upload").
		SetFormData(map[string]string{
			"name":  "John",
			"email": "john@example.com",
		})

	resp, err := client.Execute(req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Contains(t, mockClient.lastContentType, "multipart/form-data")
	assert.Contains(t, mockClient.lastContentType, "boundary=")

	bodyStr := string(capturedBody)
	assert.Contains(t, bodyStr, "John")
	assert.Contains(t, bodyStr, "john@example.com")
}

func TestMultipartFileUpload(t *testing.T) {
	t.Parallel()

	var capturedBody []byte

	mockClient := &multipartMock{}
	mockClient.onDo = func(req *http.Request) (*http.Response, error) {
		body, _ := io.ReadAll(req.Body)
		capturedBody = body
		return mockClient.mockResponse()
	}

	client := restc.NewWithClient("https://api.test.com", mockClient)

	content := "Hello, this is a test file."
	req := restc.Post("upload").
		SetFormData(map[string]string{"title": "My Document"}).
		SetFileReader("file", "test.txt", strings.NewReader(content))

	resp, err := client.Execute(req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Contains(t, mockClient.lastContentType, "multipart/form-data")

	bodyStr := string(capturedBody)
	assert.Contains(t, bodyStr, "My Document")
	assert.Contains(t, bodyStr, "test.txt")
	assert.Contains(t, bodyStr, "Hello, this is a test file.")
}

func TestMultipartFileOnly(t *testing.T) {
	t.Parallel()

	var capturedBody []byte

	mockClient := &multipartMock{}
	mockClient.onDo = func(req *http.Request) (*http.Response, error) {
		body, _ := io.ReadAll(req.Body)
		capturedBody = body
		return mockClient.mockResponse()
	}

	client := restc.NewWithClient("https://api.test.com", mockClient)

	req := restc.Post("upload").
		SetFileReader("avatar", "photo.png", bytes.NewReader([]byte{0x89, 0x50, 0x4E, 0x47}))

	resp, err := client.Execute(req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	bodyStr := string(capturedBody)
	assert.Contains(t, bodyStr, "photo.png")
}

func TestMultipartMultipleFiles(t *testing.T) {
	t.Parallel()

	var capturedBody []byte

	mockClient := &multipartMock{}
	mockClient.onDo = func(req *http.Request) (*http.Response, error) {
		body, _ := io.ReadAll(req.Body)
		capturedBody = body
		return mockClient.mockResponse()
	}

	client := restc.NewWithClient("https://api.test.com", mockClient)

	req := restc.Post("upload").
		SetFileReader("doc1", "a.txt", strings.NewReader("file A")).
		SetFileReader("doc2", "b.txt", strings.NewReader("file B"))

	resp, err := client.Execute(req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	bodyStr := string(capturedBody)
	assert.Contains(t, bodyStr, "a.txt")
	assert.Contains(t, bodyStr, "b.txt")
	assert.Contains(t, bodyStr, "file A")
	assert.Contains(t, bodyStr, "file B")
}

func TestMultipartParseMultipartForm(t *testing.T) {
	t.Parallel()

	var parsedForm *multipart.Form

	mockClient := &multipartMock{}
	mockClient.onDo = func(req *http.Request) (*http.Response, error) {
		err := req.ParseMultipartForm(10 << 20)
		if err == nil {
			parsedForm = req.MultipartForm
		}
		return mockClient.mockResponse()
	}

	client := restc.NewWithClient("https://api.test.com", mockClient)

	req := restc.Post("upload").
		SetFormData(map[string]string{"key": "value"}).
		SetFileReader("file", "test.txt", strings.NewReader("content"))

	resp, err := client.Execute(req)
	assert.NoError(t, err)
	assert.NotNil(t, parsedForm)
	assert.Equal(t, "value", parsedForm.Value["key"][0])
	assert.Equal(t, "test.txt", parsedForm.File["file"][0].Filename)
	assert.NotNil(t, resp)
}

func TestFormURLEncoded(t *testing.T) {
	t.Parallel()

	var capturedBody []byte
	var capturedContentType string

	mockClient := &multipartMock{}
	mockClient.onDo = func(req *http.Request) (*http.Response, error) {
		body, _ := io.ReadAll(req.Body)
		capturedBody = body
		capturedContentType = req.Header.Get("Content-Type")
		return mockClient.mockResponse()
	}

	client := restc.NewWithClient("https://api.test.com", mockClient)

	req := restc.Post("login").
		SetFormURLEncoded(map[string]string{
			"username": "john",
			"password": "secret",
		})

	resp, err := client.Execute(req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "application/x-www-form-urlencoded", capturedContentType)

	bodyStr := string(capturedBody)
	assert.Contains(t, bodyStr, "username=john")
	assert.Contains(t, bodyStr, "password=secret")
}
