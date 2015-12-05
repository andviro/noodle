package gorilla

import (
	"github.com/andviro/noodle"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
	"net/http"
)

type key int

var varKey key = 0

// Vars injects Gorilla mux route variables into context
func Vars(next noodle.Handler) noodle.Handler {
	return func(c context.Context, w http.ResponseWriter, r *http.Request) error {
		withVars := context.WithValue(c, varKey, mux.Vars(r))
		return next(withVars, w, r)
	}
}

// GetVars extracts route variables from context
func GetVars(c context.Context) map[string]string {
	res, _ := c.Value(varKey).(map[string]string)
	return res
}
