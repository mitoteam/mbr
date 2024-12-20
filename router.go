package mbr

import (
	"log"
	"net/http"
)

func Handler(rootController Controller) http.Handler {
	router = &mbrRouterT{
		routes: make(map[string]*Route),
	}

	for _, route := range scanRoutes(rootController) {
		route.fullPath = rootController.BasePath() + "/" + route.Path
		router.routes[route.name] = &route
	}

	mux := http.NewServeMux()

	for name, route := range router.routes {
		log.Printf("Path found %s => %s", name, route.fullPath)
		mux.HandleFunc(route.fullPath+"/{$}", buildHandleRouteFunc(route))
	}

	return mux
}

type mbrRouterT struct {
	routes map[string]*Route // path => Route
}

var router *mbrRouterT
