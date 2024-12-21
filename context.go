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
	w       http.ResponseWriter
	request *http.Request
}

func newContext(w http.ResponseWriter, r *http.Request, route *Route) *Context {
	ctx := &Context{
		//originalCtx: request.Context(), //not needed yet
		w:     w,
		route: route,
	}

	httpCtx := context.WithValue(r.Context(), mbrContextKey, ctx)
	ctx.request = r.WithContext(httpCtx)

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

func (ctx *Context) Writer() http.ResponseWriter {
	return ctx.w
}

func (ctx *Context) Request() *http.Request {
	return ctx.request
}
