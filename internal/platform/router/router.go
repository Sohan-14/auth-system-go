package router

import (
	"net/http"
	"slices"
	"strings"

	"github.com/sohan-14/auth-system-go/internal/platform/middleware"
)

// Router registers routes under a path prefix with a middleware stack.
type Router struct {
	mux    *http.ServeMux
	prefix string
	mws    []middleware.Middleware
}

// New creates a root router.
func New() *Router {
	return &Router{
		mux: http.NewServeMux(),
	}
}

// Group returns a child router. Its routes get prefix prepended and run
// through the parent's middleware first, then mws.
func (r *Router) Group(prefix string, mws ...middleware.Middleware) *Router {
	return &Router{
		mux:    r.mux,
		prefix: r.prefix + strings.TrimRight(prefix, "/"),
		mws:    append(slices.Clip(r.mws), mws...),
	}
}

// Handle registers h for a pattern like "GET /users/{id}" (method optional,
// path must start with "/"), applying this group's prefix and middleware.
func (r *Router) Handle(pattern string, h http.Handler) {
	r.mux.Handle(r.join(pattern), middleware.Chain(h, r.mws...))
}

// HandleFunc is Handle for a handler function.
func (r *Router) HandleFunc(pattern string, fn http.HandlerFunc) {
	r.Handle(pattern, fn)
}

// ServeHTTP implements http.Handler.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}

// join inserts the prefix after the optional method: "GET /x" -> "GET /v1/x".
func (r *Router) join(pattern string) string {
	if method, path, ok := strings.Cut(pattern, " "); ok {
		return method + " " + r.prefix + path
	}
	return r.prefix + pattern
}
