package middleware_test

import (
	"gopkg.in/andviro/noodle.v2"
	mw "gopkg.in/andviro/noodle.v2/middleware"
	"gopkg.in/tylerb/is.v1"
	"net/http"
	"testing"
)

func TestRealIP(t *testing.T) {
	is := is.New(t)
	n := noodle.New(mw.RealIP).
		Then(func(w http.ResponseWriter, r *http.Request) {
			realIP := mw.GetRealIP(r)
			t.Log(realIP)
			is.Equal(realIP, "testIP")
		})

	r, _ := http.NewRequest("GET", "http://localhost", nil)
	r.Header.Set("X-Real-Ip", "testIP")
	n(nil, r)

	r, _ = http.NewRequest("GET", "http://localhost", nil)
	r.Header.Set("X-Forwarded-For", "testIP, proxyIP")
	n(nil, r)
}
