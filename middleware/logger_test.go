package middleware_test

import (
	"bytes"
	"errors"
	"github.com/andviro/noodle"
	mw "github.com/andviro/noodle/middleware"
	"golang.org/x/net/context"
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
	testError := errors.New("test error")

	n := noodle.New(mw.Logger).Then(
		func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			w.WriteHeader(400)
			return testError
		},
	)

	r, _ := http.NewRequest("GET", "http://localhost", nil)
	err := n(context.TODO(), httptest.NewRecorder(), r)
	logString := buf.String()
	is.Equal(err, testError)
	is.True(strings.Contains(logString, "GET"))
	is.True(strings.Contains(logString, "(400)"))
	is.True(strings.Contains(logString, "test error"))
}
