package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

// RecoverError contains the original error value and a stack trace
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

// Recover is a basic middleware that catches panics and passes them to the
// pre-defined error handler func
func Recover(f func(error)) {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if e := recover(); e != nil {
					f(RecoverError{e, debug.Stack()})
				}
			}()
			next(w, r)
		}
	}
}
