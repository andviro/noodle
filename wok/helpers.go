package wok

import (
	"github.com/andviro/noodle"
	mw "github.com/andviro/noodle/middleware"
)

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
