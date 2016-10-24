package middleware_test

import (
	"context"
	"github.com/andviro/noodle"
	mw "github.com/andviro/noodle/middleware"
	"gopkg.in/tylerb/is.v1"
	"net/http"
	"testing"
)

func TestRealIP(t *testing.T) {
	is := is.New(t)
	n := noodle.New(mw.RealIP).
		Then(func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			realIP := mw.GetRealIP(ctx)
			t.Log(realIP)
			is.Equal(realIP, "testIP")
			return nil
		})

	r, _ := http.NewRequest("GET", "http://localhost", nil)
	r.Header.Set("X-Real-Ip", "testIP")
	is.NotErr(n(context.TODO(), nil, r))

	r, _ = http.NewRequest("GET", "http://localhost", nil)
	r.Header.Set("X-Forwarded-For", "testIP, proxyIP")
	is.NotErr(n(context.TODO(), nil, r))
}
