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
	n := noodle.New(mw.Recover(func(w http.ResponseWriter, r *http.Request, err error) {
		is.Equal(err.Error(), "panic: whoopsie!")
		_, ok := err.(mw.RecoverError)
		is.True(ok)
		http.Error(w, "fancy error code", 333)
	})).Then(panickyHandler)
	r, _ := http.NewRequest("GET", "http://localhost", nil)
	w := httptest.NewRecorder()
	n(w, r)
	is.Equal(w.Result().StatusCode, 333)
}

func panickyHandler(w http.ResponseWriter, r *http.Request) {
	panic("whoopsie!")
}
