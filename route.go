package mbr

import (
	"context"
	"fmt"
	"io/fs"
	"net/http"
	"path"
	"reflect"
	"regexp"
	"strings"

	"github.com/mitoteam/mttools"
)

type RouterHandleFunc func(ctx *MbrContext) any

type Route struct {
	signature string
	name      string
	fullPath  string
	ctrl      Controller

	PathPattern     string
	Method          string // empty = any, or space-separated methods list. examples "GET", "POST GET", "HEAD, GET"
	NotStrict       bool
	HandleF         RouterHandleFunc
	ChildController Controller
	StaticFS        fs.FS
}

func (route *Route) Name() string {
	return route.name
}

func (route *Route) FullPath() string {
	return route.fullPath
}

func (route *Route) MethodList() []string {
	s := strings.ToUpper(route.Method)

	s = regexp.MustCompile("[^A-Z]+").ReplaceAllString(s, " ")

	return strings.Fields(s)
}

// https://pkg.go.dev/net/http#ServeMux
func (route *Route) serveMuxPattern() string {
	if route.NotStrict || route.StaticFS != nil {
		return route.fullPath
	} else {
		return path.Join(route.fullPath, "{$}")
	}
}

func (route *Route) buildRouteHandler() http.Handler {
	routeHandler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if route.HandleF == nil {
			w.Write([]byte("route.HandleF is empty"))
		} else {
			//get MbrContext from request
			mbrContext := Context(r)

			//log.Println("Calling route.HandleF()")
			output := route.HandleF(mbrContext)
			//log.Println("route.HandleF() done")

			processHandlerOutput(mbrContext, w, output)
		}
	}))

	//put handler through controller's middlewares
	middlewares := route.ctrl.Middlewares()
	for i := len(middlewares) - 1; i >= 0; i-- {
		routeHandler = middlewares[i](routeHandler)
	}

	//put handler through parents middlewares
	for _, ctrl := range route.ctrl.ParentControllers() {
		middlewares := ctrl.Middlewares()
		for i := len(middlewares) - 1; i >= 0; i-- {
			routeHandler = middlewares[i](routeHandler)
		}
	}

	// add internal middleware that sets context (added last, so will be called first)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//create new context and set to request's context
		mbrContext := &MbrContext{
			//originalCtx: request.Context(), //not needed yet
			w:     w,
			route: route,
		}

		httpCtx := context.WithValue(r.Context(), mbrContextKey, mbrContext)
		r = r.WithContext(httpCtx)
		mbrContext.request = r

		//log.Println("haha! New MbrContext created")

		routeHandler.ServeHTTP(w, r)
	})
}

func processHandlerOutput(ctx *MbrContext, w http.ResponseWriter, output any) {
	switch v := output.(type) {
	case nil:
		//returning nil means "do nothing, I've done everything myself in a handler"

	case error:
		//errors issue 500 server error status
		out := fmt.Sprintf("Internal Error: %s", v.Error())
		http.Error(w, out, http.StatusInternalServerError)

	default:
		//try to convert it to string
		if v, ok := mttools.AnyToStringOk(v); ok {
			w.Write([]byte(v)) //sent string as-is
		} else {
			http.Error(w, fmt.Sprintf("Unknown handler output type: %s", reflect.TypeOf(output).String()), http.StatusInternalServerError)
		}
	}
}
