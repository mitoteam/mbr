package mbr

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"reflect"
	"regexp"
	"strings"

	"github.com/mitoteam/mttools"
)

type RouterHandleFunc func(ctx *MbrContext) any

type Route struct {
	WithMiddlewareListBase

	signature string
	name      string
	fullPath  string
	ctrl      Controller

	PathPattern     string
	Method          string // empty = any, or space-separated methods list. examples "GET", "POST GET", "HEAD, GET"
	HandleF         RouterHandleFunc
	ChildController Controller
	StaticFS        fs.FS
	FileFromFS      string
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
func (route *Route) serveMuxPattern() (pathPattern string) {
	if route.FileFromFS != "" {
		return route.fullPath
	} else if route.StaticFS != nil {
		return route.fullPath + "/"
	} else {
		return route.fullPath
	}
}

func (route *Route) buildRouteHandler() http.Handler {
	var routeHandler http.Handler
	if route.FileFromFS != "" {
		routeHandler = route.routeHandlerFileFromFS()
	} else if route.StaticFS != nil {
		routeHandler = route.routeHandlerStaticFS()
	} else {
		routeHandler = route.routeHandlerCustom()
	}

	//put handler through route's middlewares
	routeHandler = route.ApplyMiddlewares(routeHandler)

	//put handler through controller's middlewares
	routeHandler = route.ctrl.ApplyMiddlewares(routeHandler)

	//put handler through parents middlewares
	for _, parentController := range route.ctrl.ParentControllers() {
		routeHandler = parentController.ApplyMiddlewares(routeHandler)
	}

	// add internal middleware that sets context (added last, so will be called first)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//create new context and set to request's context
		mbrContext := &MbrContext{
			//originalCtx: request.Context(), //not needed yet
			w:      w,
			route:  route,
			values: mttools.NewValues(),
		}

		httpCtx := context.WithValue(r.Context(), mbrContextKey, mbrContext)
		r = r.WithContext(httpCtx)
		mbrContext.request = r

		//log.Println("DBG: New MbrContext created")

		routeHandler.ServeHTTP(w, r)
	})
}

func (route *Route) routeHandlerCustom() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	})
}

func (route *Route) routeHandlerFileFromFS() http.Handler {
	mttools.AssertNotNil(route.StaticFS, "StaticFS should be set when FileFromFS given")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		file, err := route.StaticFS.Open(route.FileFromFS)

		if err != nil {
			log.Printf("FileFromFS error: %s\nRoute: %s [%s]\n", err.Error(), route.Name(), route.FullPath())
			http.Error(w, fmt.Sprintf("Internal Error: %s", err.Error()), http.StatusInternalServerError)

			return
		}

		buf := make([]byte, 4096)
		for {
			len, err := file.Read(buf)

			if err != nil && err != io.EOF {
				log.Printf("FileFromFS error: %s\nRoute: %s [%s]\n", err.Error(), route.Name(), route.FullPath())
				http.Error(w, fmt.Sprintf("Internal Error: %s", err.Error()), http.StatusInternalServerError)

				return
			}

			if len > 0 {
				w.Write(buf)
			} else {
				break
			}
		}
	})
}

func (route *Route) routeHandlerStaticFS() http.Handler {
	staticServer := http.FileServerFS(route.StaticFS)
	return http.StripPrefix(route.serveMuxPattern(), staticServer)
}

func processHandlerOutput(ctx *MbrContext, w http.ResponseWriter, output any) {
	switch v := output.(type) {
	case nil:
		//returning nil means "do nothing, I've done everything myself in a handler"

	case error:
		//errors issue 500 server error status
		log.Printf("Error %d: %s\nRoute: %s [%s]\n", http.StatusInternalServerError, v.Error(), ctx.route.Name(), ctx.route.FullPath())
		ctx.Error("Internal server Error: %s", v.Error())

	default:
		//try to convert it to string
		if v, ok := mttools.AnyToStringOk(v); ok {
			w.Write([]byte(v)) //sent string as-is
		} else {
			ctx.Error("Unknown handler output type: %s", reflect.TypeOf(output).String())
		}
	}
}
