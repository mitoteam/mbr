package mbr

import (
	"context"
	"net/http"
)

const mbrContextKey = "mitoteam/mbrContextKey"

type Context struct {
	parent context.Context

	Route *Route
}

func newContext(ctx context.Context) *Context {
	return &Context{
		parent: ctx,
	}
}

// gets mbr.Context from request's context
func MbrContext(r *http.Request) *Context {
	if ctx, ok := r.Context().Value(mbrContextKey).(*Context); ok {
		return ctx
	}

	return nil
}
