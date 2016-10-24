package adapt_test

import (
	"context"
	"errors"
	"fmt"
	"github.com/andviro/noodle"
	"github.com/andviro/noodle/adapt"
	"gopkg.in/tylerb/is.v1"
	"net/http"
	"net/http/httptest"
	"testing"
)

func httpMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "HTTP>")
		next.ServeHTTP(w, r)
	})
}

func noodleMW(next noodle.Handler) noodle.Handler {
	return func(c context.Context, w http.ResponseWriter, r *http.Request) error {
		return next(context.WithValue(c, "testKey", "testValue"), w, r)
	}
}

func TestHttpContextPasses(t *testing.T) {
	is := is.New(t)

	n := noodle.New(noodleMW, adapt.Http(httpMW)).Then(
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

func TestHttpErrorPropagates(t *testing.T) {
	is := is.New(t)
	testError := errors.New("test error")

	n := noodle.New(noodleMW, adapt.Http(httpMW)).Then(
		func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			return testError
		},
	)
	r, _ := http.NewRequest("GET", "http://localhost", nil)
	err := n(context.TODO(), httptest.NewRecorder(), r)
	is.Err(err)
	is.Equal(err, testError)
}
