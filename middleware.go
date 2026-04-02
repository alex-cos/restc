package restc

import "sync"

type Middleware func(req *Request, next func(req *Request) (*Response, error)) (*Response, error)

type ClientMiddleware struct {
	middlewares []Middleware
	mutex       *sync.RWMutex
}

func NewClientMiddleware() *ClientMiddleware {
	return &ClientMiddleware{
		middlewares: make([]Middleware, 0),
		mutex:       &sync.RWMutex{},
	}
}

func (cm *ClientMiddleware) Use(middleware ...Middleware) {
	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	cm.middlewares = append(cm.middlewares, middleware...)
}

func (cm *ClientMiddleware) Execute(req *Request, next func(req *Request) (*Response, error)) (*Response, error) {
	cm.mutex.RLock()
	middlewares := make([]Middleware, len(cm.middlewares))
	copy(middlewares, cm.middlewares)
	cm.mutex.RUnlock()

	handler := func(r *Request) (*Response, error) {
		if len(middlewares) == 0 {
			return next(r)
		}

		current := middlewares[0]
		remaining := middlewares[1:]

		return current(r, func(req *Request) (*Response, error) {
			temp := &ClientMiddleware{
				middlewares: remaining,
				mutex:       &sync.RWMutex{},
			}
			return temp.Execute(req, next)
		})
	}

	return handler(req)
}
