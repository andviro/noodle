package gorilla_test

import (
	"gopkg.in/andviro/noodle.v2"
	"gopkg.in/andviro/noodle.v2/adapt/gorilla"
	"github.com/gorilla/mux"
	"gopkg.in/tylerb/is.v1"
	"net/http"
	"net/http/httptest"
	"testing"
)

var testKey int = 0

func noodleMW(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		next(w, noodle.WithValue(r, testKey, "testValue"))
	}
}

func TestChain(t *testing.T) {
	is := is.New(t)
	site := noodle.New(gorilla.Vars, noodleMW)

	router := mux.NewRouter()
	router.Handle("/{id}", site.Then(func(w http.ResponseWriter, r *http.Request) {
		is.Equal(gorilla.GetVars(r)["id"], "testId")
		is.Equal(noodle.Value(r, testKey).(string), "testValue")
	}))

	r, _ := http.NewRequest("GET", "http://localhost/testId", nil)
	router.ServeHTTP(httptest.NewRecorder(), r)
}
