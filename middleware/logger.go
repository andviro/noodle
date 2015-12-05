package middleware

import (
	"github.com/andviro/noodle"
	"golang.org/x/net/context"
	"log"
	"net/http"
	"time"
)

// Logger is a middleware that logs requests, along with
// request URI, handler return value and request timing
func Logger(next noodle.Handler) noodle.Handler {
	return func(c context.Context, w http.ResponseWriter, r *http.Request) (err error) {
		start := time.Now()
		err = next(c, w, r)
		end := time.Now()
		log.Printf("%s %s from %s (%s) error = %v", r.Method, r.URL.String(), r.RemoteAddr, end.Sub(start), err)
		return
	}
}
