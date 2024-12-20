package mbr

import (
	"path"
)

type Route struct {
	name      string
	fullPath  string
	Path      string
	HandleF   RouterHandleFunc
	Child     Controller
	NotStrict bool
}

func (route *Route) muxPath() string {
	if route.NotStrict {
		return route.fullPath
	} else {
		return path.Join(route.fullPath, "{$}")
	}
}