package adapt_test

import (
	"fmt"
	"github.com/andviro/noodle"
	"github.com/andviro/noodle/adapt"
	"gopkg.in/tylerb/is.v1"
	"net/http"
	"net/http/httptest"
	"testing"
)

func noodleMW(next noodle.Handler) noodle.Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		next(w, noodle.Set(r, "testKey", "testValue"))
	}
}

func negroniMW(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	fmt.Fprint(w, "HTTP>")
	next(w, r)
}

func TestNegroniContextPasses(t *testing.T) {
	is := is.New(t)

	n := noodle.New(noodleMW, adapt.Negroni(negroniMW)).Then(
		func(w http.ResponseWriter, r *http.Request) {
			val, ok := noodle.Get(r, "testKey").(string)
			is.True(ok)
			is.Equal(val, "testValue")
		},
	)
	r, _ := http.NewRequest("GET", "http://localhost", nil)
	n(httptest.NewRecorder(), r)
}
