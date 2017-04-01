package middleware_test

import (
	"gopkg.in/andviro/noodle.v2"
	mw "gopkg.in/andviro/noodle.v2/middleware"
	"gopkg.in/tylerb/is.v1"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecover(t *testing.T) {
	is := is.New(t)
	n := noodle.New(mw.Recover(func(err error) {
		is.Equal(err.Error(), "panic: whoopsie!")
		_, ok := err.(mw.RecoverError)
		is.True(ok)
	})).Then(panickyHandler)
	r, _ := http.NewRequest("GET", "http://localhost", nil)
	n(httptest.NewRecorder(), r)
}

func panickyHandler(w http.ResponseWriter, r *http.Request) {
	panic("whoopsie!")
}
