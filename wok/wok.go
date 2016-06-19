package wok

import (
	"github.com/andviro/noodle"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/context"
	"net/http"
	"path/filepath"
)

// Wok is a simple wrapper for httprouter with route groups and native support for noodle.Handler
type Wok struct {
	prefix  string
	parent  *Wok
	chain   noodle.Chain
	rootCtx context.Context
	*httprouter.Router
}

type RouteClosure func(noodle.Handler)

type key int

const (
	paramKey key = iota
)

var todoCtx = context.TODO()

// convert turns noodle.Handler into httprouter.Handle
func (wok *Wok) convert(h noodle.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		_ = h(context.WithValue(wok.context(), paramKey, p), w, r)
	}
}

// New creates new Wok initialized with middlewares.
// The resulting middleware chain will be called for all routes in Wok
func New(mws ...noodle.Middleware) *Wok {
	return &Wok{
		Router: httprouter.New(),
		chain:  noodle.New(mws...),
	}
}

// context determines context for the handler.
func (wok *Wok) context() context.Context {
	if wok.rootCtx != nil {
		return wok.rootCtx
	}
	if wok.parent != nil {
		return wok.parent.context()
	}
	return todoCtx
}

// SetRootCtx injects user-supplied context into the router.
// Note that you can set the context for subrouters.
// If subrouter context is not set explicitly, it will be inherited from its parent.
func (wok *Wok) SetContext(ctx context.Context) {
	wok.rootCtx = ctx
}

// Handle allows to attach some noodle Middlewares and a Handle to a route
func (wok *Wok) Handle(method, path string, mws ...noodle.Middleware) RouteClosure {
	chain := noodle.New(mws...)
	for router := wok; router != nil; router = router.parent {
		chain = router.chain.Use(chain...)
		path = filepath.Join(router.prefix, path)
	}
	return func(h noodle.Handler) {
		h = chain.Then(h)
		wok.Router.Handle(method, path, wok.convert(h))
	}
}

// Group starts new route group with common prefix.
// Middleware passed to Group will be used for all routes in it.
func (wok *Wok) Group(prefix string, mws ...noodle.Middleware) *Wok {
	return &Wok{
		prefix: prefix,
		parent: wok,
		Router: wok.Router,
		chain:  noodle.New(mws...),
	}
}

// Var returns route variable for context or empty string
func Var(c context.Context, name string) string {
	return c.Value(paramKey).(httprouter.Params).ByName(name)
}
