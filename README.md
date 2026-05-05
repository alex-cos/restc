# RESTC - HTTP Client in Go

[![Go Version](https://img.shields.io/badge/Go-1.25%2B-blue)](https://go.dev/)
[![Test Status](https://github.com/alex-cos/restc/actions/workflows/test.yml/badge.svg)](https://github.com/alex-cos/restc/actions/workflows/test.yml)
[![Lint Status](https://github.com/alex-cos/restc/actions/workflows/lint.yml/badge.svg)](https://github.com/alex-cos/restc/actions/workflows/lint.yml)
[![License](https://img.shields.io/badge/License-MIT-green)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/alex-cos/restc)](https://goreportcard.com/report/github.com/alex-cos/restc)

## Description

RESTC is a lightweight Go library for executing HTTP requests with support for headers, cookies, query parameters, and JSON serialization. It includes error handling, automatic retries, context support, and simplified response management.

## Features

- Supports HTTP methods: GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS, TRACE
- Easy handling of headers, cookies, and query parameters
- Automatic JSON serialization and deserialization
- Bearer token authentication
- HTTP error handling with typed error responses
- Retry mechanism with exponential backoff
- Context support for cancellation and timeouts
- Custom response/error parsers
- HTML error body text extraction
- Optional response body size limit (DoS protection)
- URL scheme validation (http/https only)
- Multipart form data and file upload
- Middleware chain for logging, tracing, metrics, etc.
- Redirect policy control (follow, block, or limit)
- Variadic options for flexible client configuration
- IPv4/IPv6 transport control

## Installation

```bash
go get github.com/alex-cos/restc
```

## Usage

### Creating an HTTP client

```go
// Basic client
client := restc.New("https://api.example.com")

// With custom timeout
client := restc.New("https://api.example.com",
    restc.WithTimeout(10 * time.Second),
)

// With custom http.Client (for TLS config, proxies, etc.)
httpClient := &http.Client{
    Transport: &http.Transport{
        MaxIdleConns: 10,
    },
}
client := restc.NewWithClient("https://api.example.com", httpClient)
```

### Configuring the client

```go
client := restc.New("https://api.example.com")

client.SetTimeout(10 * time.Second)
client.SetEntryPoint("https://api.example.com/v2")
client.SetRetryCount(3)
client.SetRetryWaitTime(100 * time.Millisecond)
client.SetRetryMaxWaitTime(2 * time.Second)
client.SetMaxResponseSize(10 * 1024 * 1024) // 10 MB limit
```

### Client options

You can pass options directly to `New` or `NewWithClient` for a more concise configuration:

```go
client := restc.New("https://api.example.com",
    restc.WithTimeout(10 * time.Second),
    restc.WithRetryCount(3),
    restc.WithRetryWaitTime(100 * time.Millisecond),
    restc.WithRetryMaxWaitTime(2 * time.Second),
    restc.WithMaxResponseSize(10 * 1024 * 1024),
    restc.WithHeader("User-Agent", "my-app/1.0"),
    restc.WithHeaders(map[string]string{
        "X-Custom-Header": "value",
    }),
    restc.WithDisableIPv6(),  // Force IPv4 only
)
```

Available options:

| Option | Description |
| -------- | ------------- |
| `WithTimeout(duration)` | Set request timeout |
| `WithRetryCount(n)` | Set number of retry attempts |
| `WithRetryWaitTime(duration)` | Set initial wait time between retries |
| `WithRetryMaxWaitTime(duration)` | Set max wait time between retries |
| `WithMaxResponseSize(bytes)` | Limit response body size (DoS protection) |
| `WithParseResponse(fn)` | Custom response parser |
| `WithParseError(fn)` | Custom error response parser |
| `WithHeader(key, value)` | Set a default header |
| `WithHeaders(map)` | Set multiple default headers |
| `WithRedirectPolicy(policy)` | Set redirect behavior |
| `WithMaxRedirects(n)` | Limit number of redirects |
| `WithDisableIPv6()` | Force IPv4-only connections |
| `WithOnlyIPv6()` | Force IPv6-only connections |

### Executing a GET request

```go
req := restc.Get("users").
    SetHeader("Accept", restc.TypeApplicationJSON).
    AddQueryParam("limit", "10")

resp, err := client.Execute(req)
if err != nil {
    log.Fatal(err)
}
fmt.Println("Status:", resp.StatusCode())
fmt.Println("Response:", resp.String())
```

### Executing a POST request

```go
type User struct {
    Name  string `json:"name"`
    Email string `json:"email"`
}

user := User{
    Name:  "John Doe",
    Email: "john@example.com",
}

req := restc.Post("users").SetBody(user)
resp, err := client.Execute(req)
if err != nil {
    log.Fatal(err)
}
```

### JSON deserialization

```go
type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

req := restc.Get("users/1").
    SetResponseType(&User{})

resp, err := client.Execute(req)
if err != nil {
    log.Fatal(err)
}

user := resp.Content().(*User)
fmt.Println(user.Name)
```

### Error response deserialization

```go
type APIError struct {
    Status  int    `json:"status"`
    Message string `json:"error"`
    Path    string `json:"path"`
}

req := restc.Get("users/1").
    SetResponseType(&User{}).
    SetErrorRespType(&APIError{})

resp, err := client.Execute(req)
if err != nil {
    log.Fatal(err)
}

if resp.IsError() {
    apiErr := resp.Content().(*APIError)
    fmt.Printf("Error %d: %s\n", apiErr.Status, apiErr.Message)
}
```

### Authentication

```go
// Bearer token (default scheme)
req := restc.Get("users").
    SetAuthToken("your-jwt-token")

// Custom auth scheme
req := restc.Get("users").
    SetAuthScheme("Basic").
    SetAuthToken(base64EncodedCredentials)
```

### Query parameters

```go
req := restc.Get("users").
    SetQueryParam("page", "1").
    SetQueryParam("limit", "20").
    SetQueryParams(map[string]string{
        "sort":  "name",
        "order": "asc",
    })
```

### Cookies

```go
req := restc.Get("users").
    SetCookie(&http.Cookie{
        Name:  "session",
        Value: "abc123",
    })
```

### Context support

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

resp, err := client.ExecuteWithContext(ctx, req)
if err != nil {
    if errors.Is(err, context.DeadlineExceeded) {
        log.Println("Request timed out")
    }
}
```

### Custom parsers

```go
client.SetParseResponse(func(request *restc.Request, response *restc.Response) (any, error) {
    // Custom parsing logic
    return myCustomParser(response.Bytes())
})

client.SetParseError(func(request *restc.Request, response *restc.Response) (any, error) {
    // Custom error parsing logic (supports HTML text extraction)
    return restc.DefaultParseError(request, response)
})
```

### Multipart form data & file upload

```go
// Form fields only
req := restc.Post("upload").
    SetFormData(map[string]string{
        "title": "My Document",
        "desc":  "A test file",
    })

// With file from disk
req := restc.Post("upload").
    SetFormData(map[string]string{"title": "Avatar"}).
    SetFile("photo", "/path/photo.jpg")

// With file from io.Reader
req := restc.Post("upload").
    SetFormData(map[string]string{"id": "123"}).
    SetFileReader("document", "report.pdf", pdfReader)

// Multiple files
req := restc.Post("upload").
    SetFileReader("doc1", "a.txt", readerA).
    SetFileReader("doc2", "b.txt", readerB)
```

### URL-encoded form data

```go
req := restc.Post("login").
    SetFormURLEncoded(map[string]string{
        "username": "john",
        "password": "secret",
    })
```

### Middleware

```go
// Logging middleware
client.UseMiddleware(func(req *restc.Request, next func(req *restc.Request) (*restc.Response, error)) (*restc.Response, error) {
    start := time.Now()
    resp, err := next(req)
    log.Printf("[%s] %s %d (%s)", req.String(), resp.Status(), resp.StatusCode(), time.Since(start))
    return resp, err
})

// Short-circuit middleware (skip execution)
client.UseMiddleware(func(req *restc.Request, next func(req *restc.Request) (*restc.Response, error)) (*restc.Response, error) {
    if req.GetAuthToken() == "" {
        return nil, errors.New("missing auth token")
    }
    return next(req)
})

// Multiple middlewares execute in order (onion model)
client.UseMiddleware(loggingMiddleware, tracingMiddleware, metricsMiddleware)
```

### Redirect policy

```go
// Block all redirects (stop at first 3xx)
client.SetRedirectPolicy(restc.NoRedirect)

// Follow redirects (default behavior)
client.SetRedirectPolicy(restc.FollowRedirects)

// Limit the number of sucessive redirects
client.SetMaxRedirects(3)
```

## Response API

```go
resp, _ := client.Execute(req)

resp.StatusCode()      // int - HTTP status code
resp.Status()          // string - "200 OK"
resp.IsError()         // bool - true if status >= 400
resp.String()          // string - response body as string
resp.Bytes()           // []byte - raw response body
resp.Content()         // any - parsed content (via SetResponseType/SetErrorRespType)
resp.Proto()           // string - "HTTP/2.0"
resp.Header()          // http.Header - response headers
resp.Cookies()         // []*http.Cookie - response cookies
resp.ContentType()     // string - Content-Type header
resp.ReceivedAt()      // time.Time - when response was received
```

## Constants

### Content types

```go
restc.TypeApplicationJSON           // "application/json"
restc.TypeApplicationXML            // "application/xml"
restc.TypeApplicationFormURLEncoded // "application/x-www-form-urlencoded"
restc.TypeMultipartFormData         // "multipart/form-data"
restc.TypeTextHTML                  // "text/html"
restc.TypeTextPLAIN                 // "text/plain"
restc.TypeTextXML                   // "text/xml"
// ... and many more
```

### HTTP headers

```go
restc.ContentType     // "Content-Type"
restc.Authorization   // "Authorization"
restc.Accept          // "Accept"
restc.UserAgent       // "User-Agent"
// ... and many more
```

### HTTP methods

```go
restc.MethodGet
restc.MethodPost
restc.MethodPut
restc.MethodPatch
restc.MethodDelete
restc.MethodHead
restc.MethodOptions
restc.MethodTrace
```

## Retry mechanism

Retries use exponential backoff with configurable wait times:

```go
client.SetRetryCount(3)                          // 3 retry attempts (4 total)
client.SetRetryWaitTime(100 * time.Millisecond)  // initial wait
client.SetRetryMaxWaitTime(2 * time.Second)      // max wait between retries
```

The retry mechanism only retries on **transient errors** (network timeouts, dial errors). Non-retriable errors (context cancellation, invalid URL scheme, parse errors) fail immediately.

## Security

- Only `http` and `https` URL schemes are accepted
- Response body size can be limited with `SetMaxResponseSize` to prevent memory exhaustion
