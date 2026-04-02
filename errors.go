package restc

import "errors"

// Client errors.
var (
	ErrMaxRedirects   = errors.New("maximum redirects exceeded")
	ErrHTTPRequest    = errors.New("failed to execute HTTP request")
	ErrReadBody       = errors.New("failed to read body response")
	ErrParseJSON      = errors.New("failed to parse JSON response")
	ErrParseXML       = errors.New("failed to parse XML response")
	ErrParseHTML      = errors.New("failed to parse HTML response")
	ErrUnexpectedType = errors.New("unexpected response type")
)

// Request errors.
var (
	ErrInvalidMethod     = errors.New("invalid HTTP method")
	ErrUnsupportedScheme = errors.New("unsupported URL scheme")
	ErrUnsupportedBody   = errors.New("unsupported body type")
	ErrInvalidEntryPoint = errors.New("failed to parse given entry point")
	ErrBuildRequest      = errors.New("failed to build HTTP request")
	ErrMultipart         = errors.New("multipart form error")
)
