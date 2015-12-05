package middleware_test

import (
	"github.com/andviro/noodle"
	mw "github.com/andviro/noodle/middleware"
	"golang.org/x/net/context"
	"gopkg.in/tylerb/is.v1"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStore(t *testing.T) {
	is := is.New(t)
	n := noodle.New(mw.LocalStore).Then(
		func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			s := mw.GetStore(ctx)
			is.NotNil(s)
			return nil
		},
	)
	r, _ := http.NewRequest("GET", "http://localhost", nil)
	_ = n(context.TODO(), httptest.NewRecorder(), r)
}
