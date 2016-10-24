package adapt

import (
	"context"
	"github.com/andviro/noodle"
	"net/http"
)

// Http converts generic "dumb" middleware to context-aware, so that context
// is maintained throgout calling chain and error value is propagated correctly
func Http(mw func(http.Handler) http.Handler) noodle.Middleware {
	return func(next noodle.Handler) noodle.Handler {
		return func(c context.Context, w http.ResponseWriter, r *http.Request) error {
			var err error
			wrappedNext := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				err = next(c, w, r)
			})
			mw(wrappedNext).ServeHTTP(w, r)
			return err
		}
	}
}
