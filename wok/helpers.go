package wok

import (
	"context"
	"github.com/andviro/noodle"
	mw "github.com/andviro/noodle/middleware"
	"net/http"
)

type HandlerC interface {
	ServeHTTPC(context.Context, http.ResponseWriter, *http.Request)
}

// HandleC adapts handlers that implement ServeHTTPC method into Wok
func (f RouteClosure) HandleC(h HandlerC) {
	f(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		h.ServeHTTPC(ctx, w, r)
		return nil
	})
}

// HandleFuncC adapts functions compatible with ServeHTTPC signature into Wok
func (f RouteClosure) HandleFuncC(h func(ctx context.Context, w http.ResponseWriter, r *http.Request)) {
	f(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		h(ctx, w, r)
		return nil
	})
}

// Default creates new Wok with the default Noodle middleware chain
func Default(mws ...noodle.Middleware) *Wok {
	return New(mw.Default(mws...)...)
}

// GET is a convenience wrapper over Wok.Handle
func (wok *Wok) GET(path string, mws ...noodle.Middleware) RouteClosure {
	return wok.Handle("GET", path, mws...)
}

// POST is a convenience wrapper over Wok.Handle
func (wok *Wok) POST(path string, mws ...noodle.Middleware) RouteClosure {
	return wok.Handle("POST", path, mws...)
}

// DELETE is a convenience wrapper over Wok.Handle
func (wok *Wok) DELETE(path string, mws ...noodle.Middleware) RouteClosure {
	return wok.Handle("DELETE", path, mws...)
}

// PATCH is a convenience wrapper over Wok.Handle
func (wok *Wok) PATCH(path string, mws ...noodle.Middleware) RouteClosure {
	return wok.Handle("PATCH", path, mws...)
}

// PUT is a convenience wrapper over Wok.Handle
func (wok *Wok) PUT(path string, mws ...noodle.Middleware) RouteClosure {
	return wok.Handle("PUT", path, mws...)
}

// OPTIONS is a convenience wrapper over Wok.Handle
func (wok *Wok) OPTIONS(path string, mws ...noodle.Middleware) RouteClosure {
	return wok.Handle("OPTIONS", path, mws...)
}
