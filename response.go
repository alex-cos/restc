package restc

import (
	"net/http"
	"strings"
	"time"
)

// Response represents an HTTP response with status, headers, and body.
// It wraps the standard library http.Response and provides additional functionality.
type Response struct {
	rawResponse *http.Response
	isRead      bool
	bodyBytes   []byte
	receivedAt  time.Time
	content     any
}

// NewResponse creates a new Response from an http.Response.
func NewResponse(rawResponse *http.Response) *Response {
	return &Response{
		rawResponse: rawResponse,
		isRead:      false,
		bodyBytes:   nil,
		receivedAt:  time.Now().UTC(),
		content:     nil,
	}
}

// String returns the response body as a trimmed string.
func (r *Response) String() string {
	return strings.TrimSpace(string(r.bodyBytes))
}

// Status returns the status string from the response (e.g., "200 OK").
func (r *Response) Status() string {
	if r.rawResponse == nil {
		return ""
	}
	return r.rawResponse.Status
}

// StatusCode returns the HTTP status code (e.g., 200, 404).
func (r *Response) StatusCode() int {
	if r.rawResponse == nil {
		return 0
	}
	return r.rawResponse.StatusCode
}

// IsError returns true if the response status code is 400 or higher.
func (r *Response) IsError() bool {
	if r.rawResponse == nil {
		return true
	}
	return r.rawResponse.StatusCode >= 400
}

// Proto returns the protocol version (e.g., "HTTP/1.1").
func (r *Response) Proto() string {
	if r.rawResponse == nil {
		return ""
	}
	return r.rawResponse.Proto
}

// Header returns the response headers.
func (r *Response) Header() http.Header {
	if r.rawResponse == nil {
		return http.Header{}
	}
	return r.rawResponse.Header
}

// Cookies returns the cookies sent in the response.
func (r *Response) Cookies() []*http.Cookie {
	if r.rawResponse == nil {
		return make([]*http.Cookie, 0)
	}
	return r.rawResponse.Cookies()
}

// ContentType returns the Content-Type header value.
func (r *Response) ContentType() string {
	if r.rawResponse == nil {
		return ""
	}
	return r.rawResponse.Header.Get(ContentType)
}

// Bytes returns the response body as bytes.
func (r *Response) Bytes() []byte {
	return r.bodyBytes
}

// ReceivedAt returns the time when the response was received.
func (r *Response) ReceivedAt() time.Time {
	return r.receivedAt
}

// Content returns the parsed response content.
// Content is available after successful parsing with DefaultParseResponse.
func (r *Response) Content() any {
	return r.content
}
