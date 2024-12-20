package mbr

import (
	"net/http"
)

func Handler(rootController Controller) http.Handler {
	router = &mbrRouterT{
		routes: make(map[string]*Route),
	}

	scanRoutes(rootController)

	router.routes["/"] = &Route{
		Path: "/",
	}

	router.routes["/aa"] = &Route{
		Path: "/aa",
	}

	mux := http.NewServeMux()

	for path, route := range router.routes {
		mux.HandleFunc(path+"/{$}", buildHandleRouteFunc(route))
	}

	return mux
}

type mbrRouterT struct {
	routes map[string]*Route // path => Route
}

var router *mbrRouterT
