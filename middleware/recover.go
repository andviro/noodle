package middleware

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"gopkg.in/andviro/noodle.v2"
)

// RecoverError contains the original error value and a stack trace
type RecoverError struct {
	Value      interface{}
	StackTrace []byte
}

func (r RecoverError) Error() string {
	return fmt.Sprintf("panic: %v", r.Value)
}

// Format implements Formatter interface
func (r RecoverError) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%+v\n", r.Value)
			fmt.Fprintf(s, "%s", string(r.StackTrace))
			return
		}
		fallthrough
	case 's':
		fmt.Fprint(s, r.Error())
	case 'q':
		fmt.Fprintf(s, "%q", r.Error())
	}
}

// Recover is a basic middleware that catches panics and passes them to the
// pre-defined error handler func
func Recover(f func(http.ResponseWriter, *http.Request, error)) noodle.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if e := recover(); e != nil {
					f(w, r, RecoverError{e, debug.Stack()})
				}
			}()
			next(w, r)
		}
	}
}
