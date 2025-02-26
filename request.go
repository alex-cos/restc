package restc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	_url "net/url"
	"reflect"
	"strings"
	"time"
)

type Request struct {
	url           string
	method        string
	authToken     string
	authScheme    string
	queryParams   _url.Values
	header        map[string]string
	cookies       []*http.Cookie
	body          any
	createdAt     time.Time
	respType      any
	errorRespType any
}

func NewRequest(method, url string) *Request {
	return &Request{
		url:           url,
		method:        method,
		authToken:     "",
		authScheme:    "Bearer",
		queryParams:   nil,
		header:        map[string]string{},
		cookies:       make([]*http.Cookie, 0),
		body:          nil,
		createdAt:     time.Now().UTC(),
		respType:      nil,
		errorRespType: nil,
	}
}

func Get(url string) *Request {
	return NewRequest(MethodGet, url)
}

func Head(url string) *Request {
	return NewRequest(MethodHead, url)
}

func Post(url string) *Request {
	return NewRequest(MethodPost, url)
}

func Put(url string) *Request {
	return NewRequest(MethodPut, url)
}

func Patch(url string) *Request {
	return NewRequest(MethodPatch, url)
}

func Delete(url string) *Request {
	return NewRequest(MethodDelete, url)
}

func Options(url string) *Request {
	return NewRequest(MethodOptions, url)
}

func Trace(url string) *Request {
	return NewRequest(MethodTrace, url)
}

func (r *Request) String() string {
	params := r.queryParams.Encode()
	str := fmt.Sprintf("[%s] /%s", r.method, strings.TrimLeft(r.url, "/"))
	if params != "" {
		str += "?" + params
	}

	return str
}

func (r *Request) SetURL(url string) *Request {
	r.url = url
	return r
}

func (r *Request) SetMethod(m string) *Request {
	r.method = m
	return r
}

func (r *Request) SetContentType(ct string) *Request {
	r.SetHeader(http.CanonicalHeaderKey(ContentType), ct)
	return r
}

func (r *Request) SetHeader(header, value string) *Request {
	r.header[header] = value
	return r
}

func (r *Request) SetHeaders(headers map[string]string) *Request {
	for h, v := range headers {
		r.SetHeader(h, v)
	}
	return r
}

func (r *Request) SetCookie(hc *http.Cookie) *Request {
	r.cookies = append(r.cookies, hc)
	return r
}

func (r *Request) SetCookies(rs []*http.Cookie) *Request {
	r.cookies = append(r.cookies, rs...)
	return r
}

func (r *Request) AddQueryParam(param, value string) *Request {
	r.ensureQueryParams()
	r.queryParams.Add(param, value)
	return r
}

func (r *Request) SetQueryParam(param, value string) *Request {
	r.ensureQueryParams()
	r.queryParams.Set(param, value)
	return r
}

func (r *Request) SetQueryParams(params map[string]string) *Request {
	for p, v := range params {
		r.SetQueryParam(p, v)
	}
	return r
}

func (r *Request) SetQueryParamsFromValues(params _url.Values) *Request {
	for p, v := range params {
		for _, pv := range v {
			r.queryParams.Add(p, pv)
		}
	}
	return r
}

func (r *Request) SetBody(body any) *Request {
	r.body = body
	return r
}

func (r *Request) SetAuthToken(authToken string) *Request {
	r.authToken = authToken
	return r
}

func (r *Request) SetAuthScheme(scheme string) *Request {
	r.authScheme = scheme
	return r
}

func (r *Request) SetResponseType(responseType any) *Request {
	r.respType = responseType
	return r
}

func (r *Request) SetErrorRespType(responseType any) *Request {
	r.errorRespType = responseType
	return r
}

func (r *Request) GetResponseType() any {
	return r.respType
}

func (r *Request) GetErrorRespType() any {
	return r.errorRespType
}

func (r *Request) computeWithContext(ctx context.Context, entryPoint string) (*http.Request, error) {
	var (
		req    *http.Request
		err    error
		reader io.Reader
	)

	if !validMethods()[r.method] {
		return nil, fmt.Errorf("invalid provided HTTP method: %s", r.method)
	}

	url, err := _url.Parse(r.url)
	if err != nil || !url.IsAbs() {
		entryPointURL, err := _url.Parse(entryPoint)
		if err != nil {
			return nil, fmt.Errorf("failed to parse given entry point: %w", err)
		}
		url = entryPointURL.JoinPath(r.url)
	}

	if r.body != nil {
		reader, err = r.toReader()
		if err != nil {
			return nil, err
		}
	}
	req, err = http.NewRequestWithContext(ctx, r.method, url.String(), reader)
	if err != nil {
		return nil, fmt.Errorf("failed to build HTTP request: %w", err)
	}
	req.URL.RawQuery = r.queryParams.Encode()

	r.applyHeaders(req)

	return req, err
}

func (r *Request) ensureQueryParams() {
	if r.queryParams == nil {
		r.queryParams = _url.Values{}
	}
}

func (r *Request) toReader() (io.Reader, error) {
	switch body := r.body.(type) {
	case string:
		return strings.NewReader(body), nil
	case []byte:
		return bytes.NewReader(body), nil
	case io.Reader:
		return body, nil
	default:
		typ := reflect.TypeOf(body)
		if typ.Kind() == reflect.Ptr {
			typ = typ.Elem()
		}
		if typ.Kind() == reflect.Struct ||
			typ.Kind() == reflect.Array ||
			typ.Kind() == reflect.Slice {
			jsonData, err := json.Marshal(body)
			if err != nil {
				return nil, err
			}
			r.SetContentType(TypeApplicationJSON)
			return bytes.NewReader(jsonData), nil
		}
		return nil, fmt.Errorf("unsupported type: %T", body)
	}
}

func (r *Request) applyHeaders(req *http.Request) {
	for key, value := range r.header {
		req.Header.Add(key, value)
	}
	if r.body != nil && req.Header.Get(ContentType) == "" {
		req.Header.Set(ContentType, TypeApplicationJSON)
	}
	if r.authToken != "" {
		req.Header.Set(Authorization, strings.TrimSpace(r.authScheme+" "+r.authToken))
	}
	for _, cookie := range r.cookies {
		req.AddCookie(cookie)
	}
}

func validMethods() map[string]bool {
	return map[string]bool{
		http.MethodGet:     true,
		http.MethodPost:    true,
		http.MethodPut:     true,
		http.MethodDelete:  true,
		http.MethodPatch:   true,
		http.MethodHead:    true,
		http.MethodOptions: true,
		http.MethodTrace:   true,
	}
}
