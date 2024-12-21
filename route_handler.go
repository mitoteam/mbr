package mbr

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/mitoteam/mttools"
)

type RouterHandleFunc func(ctx *Context) any

func buildRouteHandlerFunc(route *Route) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//check for existing mbr.Context
		mbrContext := MbrContext(r)

		//create new one and set to request's context
		if mbrContext == nil {
			mbrContext = newContext(w, r, route)
		}

		if route.HandleF == nil {
			w.Write([]byte(" route.HandleF is empty"))
		} else {
			output := route.HandleF(mbrContext)
			processHandlerOutput(mbrContext, w, output)
		}
	}
}

func processHandlerOutput(ctx *Context, w http.ResponseWriter, output any) {
	switch v := output.(type) {
	case nil:
		//returning nil means "do nothing, I've done everything myself in a handler"

	case error:
		//errors issue 500 server error status
		http.Error(w, v.Error(), http.StatusInternalServerError)

	default:
		//try to convert it to string
		if v, ok := mttools.AnyToStringOk(v); ok {
			w.Write([]byte(v)) //sent string as-is
		} else {
			http.Error(w, fmt.Sprintf("Unknown handler output type: %s", reflect.TypeOf(output).String()), http.StatusInternalServerError)
		}
	}
}
