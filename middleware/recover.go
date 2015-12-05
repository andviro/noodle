package middleware

import (
	"fmt"
	"github.com/andviro/noodle"
	"golang.org/x/net/context"
	"net/http"
)

// Recover is a basic middleware that catches panics and converts them into
// errors
func Recover(next noodle.Handler) noodle.Handler {
	return func(c context.Context, w http.ResponseWriter, r *http.Request) (err error) {
		defer func() {
			if e := recover(); e != nil {
				err = fmt.Errorf("panic: %v", e)
			}
		}()
		err = next(c, w, r)
		return
	}
}
