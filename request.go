package restc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"net/http"
	_url "net/url"
	"strings"
	"time"
)

var validMethodsMap = map[string]bool{
	http.MethodGet:     true,
	http.MethodPost:    true,
	http.MethodPut:     true,
	http.MethodDelete:  true,
	http.MethodPatch:   true,
	http.MethodHead:    true,
	http.MethodOptions: true,
	http.MethodTrace:   true,
}

type Request struct {
	url            string
	method         string
	authToken      string
	authScheme     string
	queryParams    _url.Values
	header         map[string]string
	cookies        []*http.Cookie
	body           any
	formData       map[string]string
	formURLEncoded map[string]string
	files          []*FileUpload
	multipartErr   error
	timeout        time.Duration
	createdAt      time.Time
	respType       any
	errorRespType  any
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

func (r *Request) Method() string {
	return r.method
}

func (r *Request) URL() string {
	return r.url
}

func (r *Request) Clone() *Request {
	clone := &Request{
		url:           r.url,
		method:        r.method,
		authToken:     r.authToken,
		authScheme:    r.authScheme,
		body:          r.body,
		createdAt:     r.createdAt,
		respType:      r.respType,
		errorRespType: r.errorRespType,
		timeout:       r.timeout,
		multipartErr:  r.multipartErr,
	}

	if r.queryParams != nil {
		clone.queryParams = make(_url.Values, len(r.queryParams))
		for k, v := range r.queryParams {
			clone.queryParams[k] = append([]string(nil), v...)
		}
	}

	if r.header != nil {
		clone.header = make(map[string]string, len(r.header))
		maps.Copy(clone.header, r.header)
	}

	if r.cookies != nil {
		clone.cookies = make([]*http.Cookie, len(r.cookies))
		copy(clone.cookies, r.cookies)
	}

	if r.formData != nil {
		clone.formData = make(map[string]string, len(r.formData))
		maps.Copy(clone.formData, r.formData)
	}

	if r.formURLEncoded != nil {
		clone.formURLEncoded = make(map[string]string, len(r.formURLEncoded))
		maps.Copy(clone.formURLEncoded, r.formURLEncoded)
	}

	if r.files != nil {
		clone.files = make([]*FileUpload, len(r.files))
		copy(clone.files, r.files)
	}

	return clone
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
	if r.header == nil {
		r.header = make(map[string]string)
	}
	r.header[header] = value
	return r
}

func (r *Request) SetHeaders(headers map[string]string) *Request {
	if r.header == nil {
		r.header = make(map[string]string)
	}
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
	r.ensureQueryParams()
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

func (r *Request) SetFormURLEncoded(data map[string]string) *Request {
	r.formURLEncoded = data
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

func (r *Request) SetTimeout(timeout time.Duration) *Request {
	r.timeout = timeout
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

func (r *Request) computeWithContext(
	ctx context.Context,
	entryPoint string,
	defaultHeaders map[string]string,
) (*http.Request, error) {
	var (
		req    *http.Request
		err    error
		reader io.Reader
	)

	if !validMethodsMap[r.method] {
		return nil, fmt.Errorf("%w: method='%s'", ErrInvalidMethod, r.method)
	}

	url, err := _url.Parse(r.url)
	if err != nil || !url.IsAbs() {
		entryPointURL, err := _url.Parse(entryPoint)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrInvalidEntryPoint, err)
		}
		url = entryPointURL.JoinPath(r.url)
	}

	if url.Scheme != "http" && url.Scheme != "https" {
		return nil, fmt.Errorf("%w: scheme='%s'", ErrUnsupportedScheme, url.Scheme)
	}

	if r.multipartErr != nil {
		return nil, r.multipartErr
	}

	switch {
	case len(r.formURLEncoded) > 0:
		data := _url.Values{}
		for k, v := range r.formURLEncoded {
			data.Set(k, v)
		}
		encoded := data.Encode()
		req, err = http.NewRequestWithContext(ctx, r.method, url.String(), strings.NewReader(encoded))
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrBuildRequest, err)
		}
		req.Header.Set(ContentType, TypeApplicationFormURLEncoded)
	case len(r.formData) > 0 || len(r.files) > 0:
		reader, contentType, err := r.buildMultipartBody()
		if err != nil {
			return nil, err
		}
		req, err = http.NewRequestWithContext(ctx, r.method, url.String(), reader)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrBuildRequest, err)
		}
		req.Header.Set(ContentType, contentType)
	case r.body != nil:
		reader, err = r.toReader()
		if err != nil {
			return nil, err
		}
		req, err = http.NewRequestWithContext(ctx, r.method, url.String(), reader)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrBuildRequest, err)
		}
	default:
		req, err = http.NewRequestWithContext(ctx, r.method, url.String(), nil)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrBuildRequest, err)
		}
	}

	req.URL.RawQuery = r.queryParams.Encode()

	r.applyHeaders(req, defaultHeaders)

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
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("%w type='%T': %w", ErrUnsupportedBody, body, err)
		}
		r.SetContentType(TypeApplicationJSON)
		return bytes.NewReader(jsonData), nil
	}
}

func (r *Request) applyHeaders(req *http.Request, defaultHeaders map[string]string) {
	for key, value := range defaultHeaders {
		req.Header.Set(key, value)
	}
	for key, value := range r.header {
		req.Header.Set(key, value)
	}
	if len(r.formData) == 0 &&
		len(r.files) == 0 &&
		r.body != nil &&
		req.Header.Get(ContentType) == "" {
		req.Header.Set(ContentType, TypeApplicationJSON)
	}
	if r.authToken != "" {
		req.Header.Set(Authorization, strings.TrimSpace(r.authScheme+" "+r.authToken))
	}
	for _, cookie := range r.cookies {
		req.AddCookie(cookie)
	}
}
