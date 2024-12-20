package mbr

import (
	"log"
	"net/http"
)

func buildHandleRouteFunc(route *Route) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("HandleFunc " + r.URL.Path)

		w.Write([]byte("HandleFunc " + route.Path))

		if route.HandleF == nil {
			w.Write([]byte(" route.HandleF is empty"))
		}

	}
}
