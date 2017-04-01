package adapt

import (
	"gopkg.in/andviro/noodle.v2"
	"net/http"
)

// Negroni converts function compatible with `negroni.HandlerFunc` to
// context-aware Middleware
func Negroni(mw func(http.ResponseWriter, *http.Request, http.HandlerFunc)) noodle.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			mw(w, r, http.HandlerFunc(next))
		}
	}
}
