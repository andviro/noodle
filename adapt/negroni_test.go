package adapt_test

import (
	"errors"
	"fmt"
	"github.com/andviro/noodle"
	"github.com/andviro/noodle/adapt"
	"golang.org/x/net/context"
	"gopkg.in/tylerb/is.v1"
	"net/http"
	"net/http/httptest"
	"testing"
)

func negroniMW(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	fmt.Fprint(w, "HTTP>")
	next(w, r)
}

func TestNegroniContextPasses(t *testing.T) {
	is := is.New(t)

	n := noodle.New(noodleMW, adapt.Negroni(negroniMW)).Then(
		func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			val, ok := ctx.Value("testKey").(string)
			is.True(ok)
			is.Equal(val, "testValue")
			return nil
		},
	)
	r, _ := http.NewRequest("GET", "http://localhost", nil)
	_ = n(context.TODO(), httptest.NewRecorder(), r)
}

func TestNegroniErrorPropagates(t *testing.T) {
	is := is.New(t)
	testError := errors.New("test error")

	n := noodle.New(noodleMW, adapt.Negroni(negroniMW)).Then(
		func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			return testError
		},
	)
	r, _ := http.NewRequest("GET", "http://localhost", nil)
	err := n(context.TODO(), httptest.NewRecorder(), r)
	is.Err(err)
	is.Equal(err, testError)
}
