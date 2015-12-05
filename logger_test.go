package noodle_test

import (
	"bytes"
	"github.com/andviro/noodle"
	"gopkg.in/tylerb/is.v1"
	"log"
	"strings"
	"testing"
)

func TestLogger(t *testing.T) {
	is := is.New(t)
	buf := new(bytes.Buffer)
	log.SetOutput(buf)

	h := noodle.New(noodle.Logger).Then(handlerFactory("Dummy"))
	_ = RunHTTP(h)
	log := buf.String()
	is.True(strings.Contains(log, "GET"))
}
