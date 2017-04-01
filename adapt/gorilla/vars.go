package gorilla

import (
	"gopkg.in/andviro/noodle.v2"
	"github.com/gorilla/mux"
	"net/http"
)

type key int

var varKey key = 0

// Vars injects Gorilla mux route variables into context
func Vars(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		next(w, noodle.WithValue(r, varKey, mux.Vars(r)))
	}
}

// GetVars extracts route variables from context
func GetVars(r *http.Request) map[string]string {
	res, _ := noodle.Value(r, varKey).(map[string]string)
	return res
}
