package mbr

import (
	"fmt"
	"net/http"
	"path"
	"reflect"
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
			router.mux.HandleFunc(route.muxPath(), buildRouteHandlerFunc(route))
		}
	}

	return router.mux
}

func Dump() {
	if router == nil {
		fmt.Println("mbr.Handler() never called")
	} else {
		for name, route := range router.routes {
			fmt.Printf("%s => %s\n", name, route.muxPath())
		}
	}
}

// =================== INTERNAL STUFF =======================

type mbrRouterT struct {
	routes map[string]*Route // path => Route
	mux    *http.ServeMux
}

var router *mbrRouterT

func (router *mbrRouterT) scanRoutesR(ctrl Controller, basePath string) {
	for _, route := range scanControllerMethods(ctrl) {
		route.fullPath = path.Join(basePath, route.Path)

		if route.Child != nil {
			//go deeper
			//TODO: cycle recursion check
			router.scanRoutesR(route.Child, route.fullPath)
		} else {
			router.routes[route.name] = &route
		}
	}
}

func scanControllerMethods(ctrl Controller) (routes []Route) {
	ptrType := reflect.TypeOf(ctrl)
	elementType := ptrType.Elem()

	//log.Println("scanRoutes: " + elementType.String())

	for i := 0; i < ptrType.NumMethod(); i++ {
		m := ptrType.Method(i)
		methodType := m.Type
		//log.Printf("  method %s: %+v", m.Name, methodType)

		if methodType.Kind() == reflect.Func && //it is a function
			methodType.NumIn() == 1 && methodType.In(0) == ptrType && // with one arg which is pointer receiver to struct
			methodType.NumOut() == 1 && methodType.Out(0) == reflect.TypeFor[Route]() { // returning one value and this value is Route
			//} COMMENT TO MARK if conditions end [crazy go formatting. easier to accept rather then fight]

			// call method for it to return Route struct
			route := m.Func.Call([]reflect.Value{reflect.ValueOf(ctrl)})[0].Interface().(Route)

			//give it a name from type
			route.name = elementType.String() + "." + m.Name

			routes = append(routes, route)
		}
	}

	return routes
}
