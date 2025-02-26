# RESTC - HTTP Client in Go

## Description

RESTC is a lightweight Go library for executing HTTP requests with support for headers, cookies, query parameters, and JSON serialization. It includes error handling, automatic retries, and simplified response management.

## Features

- Supports HTTP methods: GET, POST, PUT, PATCH, DELETE, etc.
- Easy handling of headers, cookies, and query parameters
- Automatic JSON serialization and deserialization
- HTTP error handling
- Retry mechanism with retryCount and retryTimeout

## Installation

With Go installed, you can install with command line interface:

```bash
go get github.com/alex-cos/restc
```

## Usage

### Creating an HTTP client

```go
client := restc.New("https://api.example.com")
client.SetTimeout(2 * time.Second)
client.SetRetryCount(3)
```

### Executing a GET request

```go
req := restc.Get("users")
  .SetHeader("Accept", "application/json")
  .AddQueryParam("limit", "10")
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
  Name: "John Doe",
  Email: "john@example.com",
}
req := restc.Post("users").SetBody(user)
resp, err := client.Execute(req)
if err != nil {
    log.Fatal(err)
}
```

### Error handling

```go
if resp.IsError() {
    fmt.Printf("Error[%d]: %s", resp.StatusCode(), resp.Status())
}
```
