package mbr

import (
	"net/http"
)

type RouterHandleFunc func(ctx *Context) any

func buildRouteHandlerFunc(route *Route) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//check for existing mbr.Context
		mbrContext := MbrContext(r)

		//create new one and set to request's context
		if mbrContext == nil {
			mbrContext = newContext(r, route)
		}

		if route.HandleF == nil {
			w.Write([]byte(" route.HandleF is empty"))
		} else {
			output := route.HandleF(mbrContext)

			if v, ok := output.(string); ok {
				w.Write([]byte(v))
			}
		}
	}
}

func processHandlerOutput() {

}
