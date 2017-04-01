package middleware_test

import (
	"gopkg.in/andviro/noodle.v2"
	mw "gopkg.in/andviro/noodle.v2/middleware"
	"gopkg.in/tylerb/is.v1"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHttpAuth(t *testing.T) {
	is := is.New(t)
	n := noodle.New(mw.HTTPAuth("test", func(u, p string) bool {
		return p == "testPassword"
	})).Then(func(w http.ResponseWriter, r *http.Request) {
		user := mw.GetUser(r)
		is.Equal(user, "testUser")
	})

	r, _ := http.NewRequest("GET", "http://localhost", nil)
	w := httptest.NewRecorder()
	n(w, r)
	is.Equal(w.Code, http.StatusUnauthorized)
	is.Equal(w.Header().Get("WWW-Authenticate"), "Basic realm=test")

	r.SetBasicAuth("testUser", "wrongPassword")
	n(w, r)

	r.SetBasicAuth("testUser", "testPassword")
	n(w, r)
}
