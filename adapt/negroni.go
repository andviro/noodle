package adapt

import (
	"context"
	"github.com/andviro/noodle"
	"net/http"
)

// Negroni converts function compatible with `negroni.HandlerFunc` to
// context-aware Middleware
func Negroni(mw func(http.ResponseWriter, *http.Request, http.HandlerFunc)) noodle.Middleware {
	return func(next noodle.Handler) noodle.Handler {
		return func(c context.Context, w http.ResponseWriter, r *http.Request) error {
			var err error
			wrappedNext := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				err = next(c, w, r)
			})
			mw(w, r, wrappedNext)
			return err
		}
	}
}
