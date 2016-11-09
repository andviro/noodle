package adapt

import (
	"github.com/andviro/noodle"
	"net/http"
)

// Negroni converts function compatible with `negroni.HandlerFunc` to
// context-aware Middleware
func Negroni(mw func(http.ResponseWriter, *http.Request, http.HandlerFunc)) noodle.Middleware {
	return func(next noodle.Handler) noodle.Handler {
		return func(w http.ResponseWriter, r *http.Request) {
			mw(w, r, http.HandlerFunc(next))
		}
	}
}
