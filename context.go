package mbr

import (
	"net/http"

	"github.com/mitoteam/mttools"
)

type mbrContextKeyType string

var mbrContextKey mbrContextKeyType = "mitoteam/mbrContextKey"

type MbrContext struct {
	//originalCtx context.Context //not needed yet

	route   *Route
	w       http.ResponseWriter
	request *http.Request

	values mttools.Values
}

// gets MbrContext from request's http.Context
func Context(r *http.Request) *MbrContext {
	//log.Println("Asking MbrContext")
	value := r.Context().Value(mbrContextKey)

	if ctx, ok := value.(*MbrContext); ok {
		//log.Println("MbrContext found")
		return ctx
	}

	return nil
}

func (ctx *MbrContext) Route() *Route {
	return ctx.route
}

func (ctx *MbrContext) Writer() http.ResponseWriter {
	return ctx.w
}

func (ctx *MbrContext) Request() *http.Request {
	return ctx.request
}

func (ctx *MbrContext) GetOk(key string) (any, bool) {
	return ctx.values.GetOk(key)
}

func (ctx *MbrContext) Get(key string) any {
	return ctx.values.Get(key)
}

func (ctx *MbrContext) Set(key string, value any) *MbrContext {
	ctx.values.Set(key, value)
	return ctx
}
