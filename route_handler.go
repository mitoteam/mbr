package mbr

import (
	"context"
	"net/http"
)

type RouterHandleFunc func(r *http.Request) any

func buildHandleRouteFunc(route *Route) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//check mbr.Context
		mbrContext := MbrContext(r)

		//create new one and set to request's context
		if mbrContext == nil {
			mbrContext = newContext(r.Context())
			mbrContext.Route = route

			ctx := context.WithValue(r.Context(), mbrContextKey, mbrContext)
			r = r.WithContext(ctx)
		}

		mbrContext.Route = route

		if route.HandleF == nil {
			w.Write([]byte(" route.HandleF is empty"))
		} else {
			output := route.HandleF(r)

			if v, ok := output.(string); ok {
				w.Write([]byte(v))
			}
		}
	}
}
