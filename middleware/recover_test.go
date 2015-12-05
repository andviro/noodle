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

func TestRecover(t *testing.T) {
	is := is.New(t)
	n := noodle.New(mw.Recover).Then(panickyHandler)
	r, _ := http.NewRequest("GET", "http://localhost", nil)
	err := n(context.TODO(), httptest.NewRecorder(), r)
	is.Equal(err.Error(), "panic: whoopsie!")
}

func panickyHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	panic("whoopsie!")
}
