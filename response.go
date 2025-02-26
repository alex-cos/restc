package restc

import (
	"net/http"
	"strings"
	"time"
)

type Response struct {
	rawResponse *http.Response
	isRead      bool
	bodyBytes   []byte
	receivedAt  time.Time
	content     any
}

func NewResponse(rawResponse *http.Response) *Response {
	return &Response{
		rawResponse: rawResponse,
		isRead:      false,
		bodyBytes:   nil,
		receivedAt:  time.Now().UTC(),
		content:     nil,
	}
}

func (r *Response) String() string {
	return strings.TrimSpace(string(r.bodyBytes))
}

func (r *Response) Status() string {
	if r.rawResponse == nil {
		return ""
	}
	return r.rawResponse.Status
}

func (r *Response) StatusCode() int {
	if r.rawResponse == nil {
		return 0
	}
	return r.rawResponse.StatusCode
}

func (r *Response) IsError() bool {
	if r.rawResponse == nil {
		return true
	}
	return r.rawResponse.StatusCode >= 400
}

func (r *Response) Proto() string {
	if r.rawResponse == nil {
		return ""
	}
	return r.rawResponse.Proto
}

func (r *Response) Header() http.Header {
	if r.rawResponse == nil {
		return http.Header{}
	}
	return r.rawResponse.Header
}

func (r *Response) Cookies() []*http.Cookie {
	if r.rawResponse == nil {
		return make([]*http.Cookie, 0)
	}
	return r.rawResponse.Cookies()
}

func (r *Response) ContentType() string {
	if r.rawResponse == nil {
		return ""
	}
	return r.rawResponse.Header.Get(ContentType)
}

func (r *Response) Bytes() []byte {
	return r.bodyBytes
}

func (r *Response) ReceivedAt() time.Time {
	return r.receivedAt
}

func (r *Response) Content() any {
	return r.content
}
