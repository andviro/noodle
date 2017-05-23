package middleware_test

import (
	mw "gopkg.in/andviro/noodle.v2/middleware"
	"gopkg.in/tylerb/is.v1"

	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDefault(t *testing.T) {
	is := is.New(t)
	def := mw.Default()
	is.NotNil(def)
}

func TestDefaultLogsPanic(t *testing.T) {
	is := is.New(t)
	buf := new(bytes.Buffer)
	log.SetOutput(buf)

	n := mw.Default().Then(
		func(w http.ResponseWriter, r *http.Request) {
			panic("ahaahhahahah")
		},
	)

	r, _ := http.NewRequest("GET", "http://localhost", nil)
	n(httptest.NewRecorder(), r)
	logString := buf.String()
	is.True(strings.Contains(logString, "GET"))
	is.True(strings.Contains(logString, "(500)"))
}
