package mbr

import (
	"net/http"
)

// standard net/http middleware alias
type Middleware func(next http.Handler) http.Handler

type MiddlewareList []Middleware

type WithMiddlewareList interface {
	With(middleware Middleware)
	Middlewares() MiddlewareList
	ApplyMiddlewares(http.Handler) http.Handler
}

type WithMiddlewareListBase struct {
	middlewares MiddlewareList
}

func (o *WithMiddlewareListBase) With(middleware Middleware) {
	o.middlewares = append(o.middlewares, middleware)
}

func (o *WithMiddlewareListBase) Middlewares() MiddlewareList {
	return o.middlewares
}

func (o *WithMiddlewareListBase) ApplyMiddlewares(handler http.Handler) http.Handler {
	for i := len(o.middlewares) - 1; i >= 0; i-- {
		handler = o.middlewares[i](handler)
	}

	return handler
}
