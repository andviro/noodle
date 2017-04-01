package middleware_test

import (
	"bytes"
	"gopkg.in/andviro/noodle.v2"
	mw "gopkg.in/andviro/noodle.v2/middleware"
	"gopkg.in/tylerb/is.v1"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLogger(t *testing.T) {
	is := is.New(t)
	buf := new(bytes.Buffer)
	log.SetOutput(buf)

	n := noodle.New(mw.Logger).Then(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(400)
		},
	)

	r, _ := http.NewRequest("GET", "http://localhost", nil)
	n(httptest.NewRecorder(), r)
	logString := buf.String()
	is.True(strings.Contains(logString, "GET"))
	is.True(strings.Contains(logString, "(400)"))
}

func TestLoggerImplementsInterfaces(t *testing.T) {
	is := is.New(t)

	n := noodle.New(mw.Logger).Then(
		func(w http.ResponseWriter, r *http.Request) {
			_, ok := w.(http.Flusher)
			is.True(ok)
			_, ok = w.(http.Hijacker)
			is.True(ok)
			_, ok = w.(http.CloseNotifier)
			is.True(ok)
		},
	)

	r, _ := http.NewRequest("GET", "http://localhost", nil)
	n(httptest.NewRecorder(), r)
}
