package mbr

import (
	"context"
	"net/http"
)

type mbrContextKeyType string

var mbrContextKey mbrContextKeyType = "mitoteam/mbrContextKey"

type Context struct {
	//originalCtx context.Context //not needed yet

	route   *Route
	request *http.Request
}

func newContext(request *http.Request, route *Route) *Context {
	ctx := &Context{
		//originalCtx: request.Context(), //not needed yet
		route: route,
	}

	httpCtx := context.WithValue(request.Context(), mbrContextKey, ctx)
	ctx.request = request.WithContext(httpCtx)

	return ctx
}

// gets mbr.Context from request's http.Context
func MbrContext(r *http.Request) *Context {
	if ctx, ok := r.Context().Value(mbrContextKey).(*Context); ok {
		return ctx
	}

	return nil
}

func (ctx *Context) Route() *Route {
	return ctx.route
}

func (ctx *Context) Request() *http.Request {
	return ctx.request
}
