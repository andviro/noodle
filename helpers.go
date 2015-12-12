package wok

import (
	"github.com/andviro/noodle"
	mw "github.com/andviro/noodle/middleware"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/context"
)

// Default creates new Wok with Logger, Recovery and LocalStore middleware at the start of middleware chain
func Default(mws ...noodle.Middleware) *Wok {
	return New(mw.Default(mws...)...)
}

// Var returns route variable for context or empty string
func Var(c context.Context, name string) string {
	return c.Value(paramKey).(httprouter.Params).ByName(name)
}

// GET is a convenience wrapper over Wok.Handle
func (wok *Wok) GET(path string, h noodle.Handler) {
	wok.Mount("GET", path, h)
}

// POST is a convenience wrapper over Wok.Handle
func (wok *Wok) POST(path string, h noodle.Handler) {
	wok.Mount("POST", path, h)
}

// DELETE is a convenience wrapper over Wok.Handle
func (wok *Wok) DELETE(path string, h noodle.Handler) {
	wok.Mount("DELETE", path, h)
}

// PATCH is a convenience wrapper over Wok.Handle
func (wok *Wok) PATCH(path string, h noodle.Handler) {
	wok.Mount("PATCH", path, h)
}

// PUT is a convenience wrapper over Wok.Handle
func (wok *Wok) PUT(path string, h noodle.Handler) {
	wok.Mount("PUT", path, h)
}

// OPTIONS is a convenience wrapper over Wok.Handle
func (wok *Wok) OPTIONS(path string, h noodle.Handler) {
	wok.Mount("OPTIONS", path, h)
}
