package wok

import (
	"github.com/andviro/noodle"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/context"
	"net/http"
	"strings"
)

type key int

// Wok is a simple wrapper for httprouter with route groups and native support for noodle.Handler
type Wok struct {
	prefix  string
	chain   noodle.Chain
	rootCtx context.Context
	*httprouter.Router
}

var methods = []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"}
var paramKey key = 0

func restorePrefix(prefix string) noodle.Middleware {
	return func(next noodle.Handler) noodle.Handler {
		if prefix == "" {
			return next
		}
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			r.URL.Path = prefix + r.URL.Path
			return next(ctx, w, r)
		}
	}
}

// convert turns noodle.Handler into httprouter.Handle
func (wok *Wok) convert(h noodle.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		_ = h(context.WithValue(wok.rootCtx, paramKey, p), w, r)
	}
}

// ServeHTTP allows Wok to implement http.Handler interface
func (wok *Wok) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Prefix needs to be stripped do maintain composeability
	http.StripPrefix(wok.prefix, wok.Router).ServeHTTP(w, r)
}

// New creates new Wok initialized with middlewares.
// The resulting middleware chain will be called for all routes in Wok
func New(mws ...noodle.Middleware) *Wok {
	return &Wok{
		Router:  httprouter.New(),
		chain:   noodle.New(mws...),
		rootCtx: context.TODO(),
	}
}

// Mount allows to attach noodle.Handler to a route
func (wok *Wok) Mount(method, path string, h noodle.Handler) {
	wok.Handle(method, path, wok.convert(wok.chain.Then(h)))
}

// Group starts new route group with common prefix.
// Middleware passed to Group will be used for all routes in it.
func (wok *Wok) Group(prefix string, mws ...noodle.Middleware) *Wok {
	if strings.ContainsAny(prefix, ":*") {
		panic("Group prefix should not have parameters")
	}

	res := &Wok{
		prefix: prefix,
		Router: httprouter.New(),
		// prefix needs to be stuffed back into request path to fool middlewares
		chain:   noodle.New(restorePrefix(prefix)).Use(wok.chain...).Use(mws...),
		rootCtx: wok.rootCtx,
	}

	res.prefix = prefix
	for _, m := range methods {
		wok.Handler(m, prefix+"/*rest", res)
	}
	return res
}
