package gorilla_test

import (
	"github.com/andviro/noodle"
	"github.com/andviro/noodle/adapt/gorilla"
	"github.com/gorilla/mux"
	"gopkg.in/tylerb/is.v1"
	"net/http"
	"net/http/httptest"
	"testing"
)

var testKey int = 0

func noodleMW(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		next(w, noodle.Set(r, testKey, "testValue"))
	}
}

func TestChain(t *testing.T) {
	is := is.New(t)
	site := noodle.New(gorilla.Vars, noodleMW)

	router := mux.NewRouter()
	router.Handle("/{id}", site.Then(func(w http.ResponseWriter, r *http.Request) {
		is.Equal(gorilla.GetVars(r)["id"], "testId")
		is.Equal(noodle.Get(r, testKey).(string), "testValue")
	}))

	r, _ := http.NewRequest("GET", "http://localhost/testId", nil)
	router.ServeHTTP(httptest.NewRecorder(), r)
}
