package restc_test

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/alex-cos/restc"
	"github.com/stretchr/testify/mock"
)

var httpStatusCodes = map[int]string{
	http.StatusOK:                  "OK",
	http.StatusCreated:             "Created",
	http.StatusAccepted:            "Accepted",
	http.StatusNoContent:           "No Content",
	http.StatusMovedPermanently:    "Moved Permanently",
	http.StatusFound:               "Found",
	http.StatusNotModified:         "Not Modified",
	http.StatusBadRequest:          "Bad Request",
	http.StatusUnauthorized:        "Unauthorized",
	http.StatusForbidden:           "Forbidden",
	http.StatusNotFound:            "Not Found",
	http.StatusMethodNotAllowed:    "Method Not Allowed",
	http.StatusInternalServerError: "Internal Server Error",
	http.StatusNotImplemented:      "Not Implemented",
	http.StatusBadGateway:          "Bad Gateway",
	http.StatusServiceUnavailable:  "Service Unavailable",
}

type Mock struct {
	mock.Mock

	Code    int
	Content string
	Header  http.Header
}

func NewMockClient(code int, content string) restc.HTTPClient {
	mockClient := &Mock{
		Code:    code,
		Content: content,
		Header:  http.Header{},
	}
	mockClient.create()

	return mockClient
}

func (m *Mock) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)

	resp, ok := args.Get(0).(*http.Response)
	if !ok {
		return nil, errors.New("wrong type")
	}
	resp.Request = req
	if req.Method == "HEAD" {
		resp.Body = nil
	}

	return resp, args.Error(1)
}

func (m *Mock) SetContent(content string) {
	m.Content = content
	m.create()
}

func (m *Mock) SetContentType(ct string) {
	m.Header.Set(restc.ContentType, ct)
	m.create()
}

func (m *Mock) AddHeader(key, value string) {
	m.Header.Add(key, value)
	m.create()
}

func (m *Mock) create() {
	var status string

	if m.Header.Get(restc.ContentType) == "" {
		m.Header.Set(restc.ContentType, "application/json; charset=utf-8")
	}

	body := io.NopCloser(strings.NewReader(m.Content))

	if message, exists := httpStatusCodes[m.Code]; exists {
		status = fmt.Sprintf("%d %s", m.Code, message)
	} else {
		status = "Unknown HTTP code"
	}

	m.On("Do", mock.Anything).Unset()
	m.On("Do", mock.Anything).Return(&http.Response{
		Status:           status,
		StatusCode:       m.Code,
		Proto:            "HTTP/2.0",
		ProtoMajor:       2,
		ProtoMinor:       0,
		Header:           m.Header,
		Body:             body,
		ContentLength:    int64(len(m.Content)),
		TransferEncoding: []string{},
		Close:            false,
		Uncompressed:     false,
		Trailer:          http.Header{},
		Request:          nil,
		TLS:              nil,
	}, nil)
}
