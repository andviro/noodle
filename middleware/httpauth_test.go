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

func TestHttpAuth(t *testing.T) {
	is := is.New(t)
	n := noodle.New(mw.HTTPAuth("test", func(u, p string) bool {
		return p == "testPassword"
	})).Then(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		user := mw.GetUser(ctx)
		is.Equal(user, "testUser")
		return nil
	})

	r, _ := http.NewRequest("GET", "http://localhost", nil)
	w := httptest.NewRecorder()
	err := n(context.TODO(), w, r)
	is.Err(err)
	is.Equal(err, mw.UnauthorizedRequest)
	is.Equal(w.Code, http.StatusUnauthorized)

	r.SetBasicAuth("testUser", "wrongPassword")
	is.Err(n(context.TODO(), w, r))

	r.SetBasicAuth("testUser", "testPassword")
	is.NotErr(n(context.TODO(), w, r))
}
