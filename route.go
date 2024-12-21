package mbr

import (
	"path"
	"regexp"
	"strings"
)

type Route struct {
	name      string
	fullPath  string
	Pattern   string
	Method    string // empty = any, or space-separated methods list. examples "GET", "POST GET", "HEAD, GET"
	NotStrict bool

	HandleF RouterHandleFunc
	Child   Controller
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
	if route.NotStrict {
		return route.fullPath
	} else {
		return path.Join(route.fullPath, "{$}")
	}
}
