package wok

import (
	"github.com/andviro/noodle"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/context"
	"net/http"
	"path/filepath"
)

type key int

// Wok is a simple wrapper for httprouter with route groups and native support for noodle.Handler
type Wok struct {
	prefix  string
	parent  *Wok
	chain   noodle.Chain
	rootCtx context.Context
	*httprouter.Router
}

var paramKey key = 0

// convert turns noodle.Handler into httprouter.Handle
func (wok *Wok) convert(h noodle.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		_ = h(context.WithValue(wok.rootCtx, paramKey, p), w, r)
	}
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

// Handle allows to attach some noodle Middlewares and a Handle to a route
func (wok *Wok) Handle(method, path string, mws ...noodle.Middleware) func(noodle.Handler) {
	chain := wok.chain.Use(mws...)
	if wok.parent == nil {
		return func(h noodle.Handler) {
			h = chain.Then(h)
			wok.Router.Handle(method, path, wok.convert(h))
		}
	}
	return wok.parent.Handle(method, filepath.Join(wok.prefix, path), chain...)
}

// Group starts new route group with common prefix.
// Middleware passed to Group will be used for all routes in it.
func (wok *Wok) Group(prefix string, mws ...noodle.Middleware) *Wok {
	return &Wok{
		prefix:  prefix,
		parent:  wok,
		Router:  wok.Router,
		chain:   noodle.New(mws...),
		rootCtx: wok.rootCtx,
	}
}

// Var returns route variable for context or empty string
func Var(c context.Context, name string) string {
	return c.Value(paramKey).(httprouter.Params).ByName(name)
}
