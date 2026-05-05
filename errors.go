package restc

import "errors"

// Client errors.
var (
	// ErrMaxRedirects is returned when the maximum number of redirects is exceeded.
	ErrMaxRedirects = errors.New("maximum redirects exceeded")
	// ErrHTTPRequest is returned when an HTTP request fails.
	ErrHTTPRequest = errors.New("failed to execute HTTP request")
	// ErrReadBody is returned when reading the response body fails.
	ErrReadBody = errors.New("failed to read body response")
	// ErrParseJSON is returned when parsing a JSON response fails.
	ErrParseJSON = errors.New("failed to parse JSON response")
	// ErrParseXML is returned when parsing an XML response fails.
	ErrParseXML = errors.New("failed to parse XML response")
	// ErrParseHTML is returned when parsing an HTML response fails.
	ErrParseHTML = errors.New("failed to parse HTML response")
	// ErrUnexpectedType is returned when the response content type is not supported.
	ErrUnexpectedType = errors.New("unexpected response type")
)

// Request errors.
var (
	// ErrInvalidMethod is returned when an invalid HTTP method is used.
	ErrInvalidMethod = errors.New("invalid HTTP method")
	// ErrUnsupportedScheme is returned when an unsupported URL scheme is used.
	ErrUnsupportedScheme = errors.New("unsupported URL scheme")
	// ErrUnsupportedBody is returned when an unsupported body type is used.
	ErrUnsupportedBody = errors.New("unsupported body type")
	// ErrInvalidEntryPoint is returned when the entry point URL is invalid.
	ErrInvalidEntryPoint = errors.New("failed to parse given entry point")
	// ErrBuildRequest is returned when building the HTTP request fails.
	ErrBuildRequest = errors.New("failed to build HTTP request")
	// ErrMultipart is returned when a multipart form error occurs.
	ErrMultipart = errors.New("multipart form error")
)
