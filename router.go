package mbr

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"reflect"
	"runtime"
	"slices"
	"strings"

	"github.com/mitoteam/mttools"
)

func Handler(rootController Controller) http.Handler {
	if router == nil {
		router = &mbrRouterT{
			routes: make(map[string]*Route),
			mux:    http.NewServeMux(),
		}

		basePath := "/" //we are at the very root controller
		router.scanRoutesR(rootController, basePath)

		for _, route := range router.routes {
			//log.Printf("Route: %s => %s", name, route.muxPath())

			if route.Method == "" {
				//register for any method (no method specified in pattern)
				router.mux.Handle(route.serveMuxPattern(), route.buildRouteHandler())
			} else {
				//register for each specified method
				for _, method := range route.MethodList() {
					router.mux.Handle(method+" "+route.serveMuxPattern(), route.buildRouteHandler())
				}
			}
		}
	}

	return router.mux
}

func Dump() {
	if router == nil {
		panic("mbr.Handler() never called")
	} else {
		for _, route := range router.routes {
			var sb strings.Builder

			if route.signature != "" {
				sb.WriteString("[" + route.signature + "] ")
			}

			sb.WriteString(route.Name() + " => ")

			methods := strings.Join(route.MethodList(), " ")
			if methods != "" {
				sb.WriteString(methods + ": ")
			}

			sb.WriteString(route.serveMuxPattern())

			fmt.Println(sb.String())
		}
	}
}

func UrlE(routeRef any, args ...any) (r string, err error) {
	if router == nil {
		return "", errors.New("router is not initialized")
	}

	funcT := reflect.TypeOf(routeRef)
	funcV := reflect.ValueOf(routeRef)

	if funcT.Kind() != reflect.Func {
		return "", errors.New("f is not a func")
	}

	mSignature, ok := routeMethodSignature(funcV)
	if !ok {
		return "", errors.New("can not calculate method's signature")
	}

	route, ok := router.routes[mSignature]
	if !ok {
		return "", fmt.Errorf("Unknown route: %s", mSignature)
	}

	if len(args)%2 != 0 {
		return "", fmt.Errorf("args count should be even")
	}

	r = route.FullPath()
	queryValues := url.Values{}

	for i := 0; i < len(args); i += 2 {
		//log.Printf("Arg %s: %v\n", args[i], args[i+1])
		argName := mttools.AnyToString(args[i])
		argValue := mttools.AnyToString(args[i+1])

		if argName != "" {
			if strings.Contains(r, "{"+argName+"}") {
				if argValue == "" {
					return "", fmt.Errorf("Path values can not be empty (empty value for '%s')", argName)
				}

				r = strings.ReplaceAll(r, "{"+argName+"}", url.QueryEscape(argValue))
			} else {
				queryValues.Set(argName, argValue)
			}
		}
	}

	if len(queryValues) > 0 {
		r += "?" + queryValues.Encode()
	}

	return r, nil
}

func Url(routeRef any, args ...any) (r string) {
	url, err := UrlE(routeRef, args...)

	if err != nil {
		panic(err)
	}

	return url
}

// =================== INTERNAL STUFF =======================

type mbrRouterT struct {
	routes map[string]*Route // path => Route
	mux    *http.ServeMux
}

var router *mbrRouterT

func (router *mbrRouterT) scanRoutesR(ctrl Controller, basePath string) {
	for _, route := range scanControllerMethods(ctrl) {
		route.fullPath = path.Join(basePath, route.PathPattern)

		if route.ChildController != nil {
			//go deeper
			//TODO: cycle recursion check

			parents := slices.Clone(ctrl.ParentControllers())
			parents = append(parents, ctrl)
			route.ChildController.SetParentControllers(parents)

			router.scanRoutesR(route.ChildController, route.fullPath)
		} else {
			router.routes[route.signature] = &route
		}
	}
}

func scanControllerMethods(ctrl Controller) (routes []Route) {
	ctrlPointerType := reflect.TypeOf(ctrl)
	ctrlElementType := ctrlPointerType.Elem()

	//log.Println("scanRoutes: " + elementType.String())

	for i := 0; i < ctrlPointerType.NumMethod(); i++ {
		m := ctrlPointerType.Method(i)
		methodType := m.Type
		//log.Printf("  method %s: %+v", m.Name, methodType)

		if methodType.Kind() == reflect.Func && //it is a function
			methodType.NumIn() == 1 && methodType.In(0) == ctrlPointerType && // with one arg which is pointer receiver to struct
			methodType.NumOut() == 1 && methodType.Out(0) == reflect.TypeFor[Route]() { // returning one value and this value is Route
			//} COMMENT TO MARK if conditions end [crazy go formatting. easier to accept rather then fight]

			// call method for it to return Route struct
			route := m.Func.Call([]reflect.Value{reflect.ValueOf(ctrl)})[0].Interface().(Route)

			//route.dbg = fmt.Sprintf("%#v", reflect.New(methodType))
			//route.dbg = runtime.FuncForPC(m.Func.Pointer()).Name()
			route.signature, _ = routeMethodSignature(m.Func)

			//give it a name from type
			route.name = ctrlElementType.String() + "." + m.Name
			route.ctrl = ctrl

			routes = append(routes, route)
		}
	}

	return routes
}

func routeMethodSignature(funcPointerValue reflect.Value) (string, bool) {
	rf := runtime.FuncForPC(funcPointerValue.Pointer())

	if rf == nil {
		return "", false
	}

	r := rf.Name()

	if r == "" {
		return "", false
	}

	//TODO: Dig deeper for this "-fm" suffix for indirect method calls
	// https://docs.google.com/document/d/1bMwCey-gmqZVTpRax-ESeVuZGmjwbocYs1iHplK-cjo/pub
	// https://www.reddit.com/r/golang/comments/1hkq71e/question_reflection_of_struct_method_pointer/
	r, _ = strings.CutSuffix(r, "-fm")

	return r, true
}
