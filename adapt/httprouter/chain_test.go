package httprouter_test

import (
	"github.com/andviro/noodle"
	hr "github.com/andviro/noodle/adapt/httprouter"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/context"
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
	site := hr.Default(noodleMW)

	router := httprouter.New()
	router.GET("/:id", site.Then(func(c context.Context, w http.ResponseWriter, r *http.Request) error {
		is.Equal(hr.GetParams(c).ByName("id"), "testId")
		is.Equal(c.Value(testKey).(string), "testValue")
		return nil
	}))

	r, _ := http.NewRequest("GET", "http://localhost/testId", nil)
	router.ServeHTTP(httptest.NewRecorder(), r)
}
