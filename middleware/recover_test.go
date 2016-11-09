package middleware_test

import (
	"github.com/andviro/noodle"
	mw "github.com/andviro/noodle/middleware"
	//"gopkg.in/tylerb/is.v1"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRecover(t *testing.T) {
	//is := is.New(t)
	n := noodle.New(mw.Recover).Then(panickyHandler)
	r, _ := http.NewRequest("GET", "http://localhost", nil)
	n(httptest.NewRecorder(), r)
	//is.Equal(err.Error(), "panic: whoopsie!")
	//_, ok := err.(mw.RecoverError)
	//is.True(ok)
}

func panickyHandler(w http.ResponseWriter, r *http.Request) {
	panic("whoopsie!")
}
