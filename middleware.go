package mbr

import (
	"net/http"
)

// standard net/http middleware alias
type Middleware func(next http.Handler) http.Handler

type MiddlewareList []Middleware
