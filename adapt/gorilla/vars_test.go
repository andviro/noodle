package gorilla_test

import (
	"context"
	"github.com/andviro/noodle"
	"github.com/andviro/noodle/adapt/gorilla"
	"github.com/gorilla/mux"
	"gopkg.in/tylerb/is.v1"
	"net/http"
	"net/http/httptest"
	"testing"
)

var testKey int = 0

func noodleMW(next noodle.Handler) noodle.Handler {
	return func(c context.Context, w http.ResponseWriter, r *http.Request) error {
		return next(context.WithValue(c, testKey, "testValue"), w, r)
	}
}

func TestChain(t *testing.T) {
	is := is.New(t)
	site := noodle.New(gorilla.Vars, noodleMW)

	router := mux.NewRouter()
	router.Handle("/{id}", site.Then(func(c context.Context, w http.ResponseWriter, r *http.Request) error {
		is.Equal(gorilla.GetVars(c)["id"], "testId")
		is.Equal(c.Value(testKey).(string), "testValue")
		return nil
	}))

	r, _ := http.NewRequest("GET", "http://localhost/testId", nil)
	router.ServeHTTP(httptest.NewRecorder(), r)
}
