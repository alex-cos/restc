package restc

import "sync"

// HandlerFunc is a function that handles an HTTP request and returns a response or error.
type HandlerFunc func(req *Request) (*Response, error)

// Middleware is a function that intercepts the request/response flow.
// It can modify the request before passing it to the next handler,
// or modify the response after receiving it.
type Middleware func(req *Request, next HandlerFunc) (*Response, error)

// ClientMiddleware manages a chain of middleware to be executed for each request.
type ClientMiddleware struct {
	middlewares []Middleware
	mutex       *sync.RWMutex
}

// NewClientMiddleware creates a new empty ClientMiddleware.
func NewClientMiddleware() *ClientMiddleware {
	return &ClientMiddleware{
		middlewares: make([]Middleware, 0),
		mutex:       &sync.RWMutex{},
	}
}

// Use adds middleware to the middleware chain.
func (cm *ClientMiddleware) Use(middleware ...Middleware) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	cm.middlewares = append(cm.middlewares, middleware...)
}

// Execute runs the middleware chain with the given request and final handler.
func (cm *ClientMiddleware) Execute(req *Request, next HandlerFunc) (*Response, error) {
	cm.mutex.RLock()
	middlewares := make([]Middleware, len(cm.middlewares))
	copy(middlewares, cm.middlewares)
	cm.mutex.RUnlock()

	handler := next
	for i := len(middlewares) - 1; i >= 0; i-- {
		m := middlewares[i]
		nextHandler := handler
		handler = func(r *Request) (*Response, error) {
			return m(r, nextHandler)
		}
	}

	return handler(req)
}
