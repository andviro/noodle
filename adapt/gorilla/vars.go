package gorilla

import (
	"github.com/andviro/noodle"
	"github.com/gorilla/mux"
	"net/http"
)

type key int

var varKey key = 0

// Vars injects Gorilla mux route variables into context
func Vars(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		next(w, noodle.Set(r, varKey, mux.Vars(r)))
	}
}

// GetVars extracts route variables from context
func GetVars(r *http.Request) map[string]string {
	res, _ := noodle.Get(r, varKey).(map[string]string)
	return res
}
