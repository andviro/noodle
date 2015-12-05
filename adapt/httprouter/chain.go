package httprouter

import (
	"github.com/andviro/noodle"
	mw "github.com/andviro/noodle/middleware"
	hr "github.com/julienschmidt/httprouter"
	"golang.org/x/net/context"
	"net/http"
)

type key int

var paramKey key = 0

// Chain mimics noodle.Chain functionality
type Chain struct {
	noodle.Chain
}

// Then works like noodle.Chain Then method but returns httprouter-compatible Handler
// with Params injected into Context
func (a Chain) Then(h noodle.Handler) hr.Handle {
	h = a.Chain.Then(h)
	return func(w http.ResponseWriter, r *http.Request, p hr.Params) {
		ctx := context.WithValue(context.Background(), paramKey, p)
		h(ctx, w, r)
	}
}

// GetParams extracts httprouter.Params from context
func GetParams(c context.Context) hr.Params {
	res, _ := c.Value(paramKey).(hr.Params)
	return res
}

// Use mimics noodle.Chain method
func (a Chain) Use(m ...noodle.Middleware) Chain {
	return Chain{a.Chain.Use(m...)}
}

// New returns new httprouter Chain adapter
func New(m ...noodle.Middleware) Chain {
	return Chain{noodle.New(m...)}
}

// New returns new httprouter Chain adapter with pre-installed middlewares
func Default(m ...noodle.Middleware) Chain {
	return Chain{mw.Default(m...)}
}
