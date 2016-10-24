package middleware

import (
	"context"
	"fmt"
	"github.com/andviro/noodle"
	"net/http"
	"runtime/debug"
)

type RecoverError struct {
	Value      interface{}
	StackTrace []byte
}

func (r RecoverError) Error() string {
	return fmt.Sprintf("panic: %v", r.Value)
}

func (r RecoverError) String() string {
	return fmt.Sprintf("%v\n%s", r.Value, string(r.StackTrace))
}

// Recover is a basic middleware that catches panics and converts them into
// errors
func Recover(next noodle.Handler) noodle.Handler {
	return func(c context.Context, w http.ResponseWriter, r *http.Request) (err error) {
		defer func() {
			if e := recover(); e != nil {
				err = RecoverError{e, debug.Stack()}
			}
		}()
		err = next(c, w, r)
		return
	}
}
