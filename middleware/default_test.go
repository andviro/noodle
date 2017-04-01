package middleware_test

import (
	mw "gopkg.in/andviro/noodle.v2/middleware"
	"gopkg.in/tylerb/is.v1"
	"testing"
)

func TestDefault(t *testing.T) {
	is := is.New(t)
	def := mw.Default()
	is.NotNil(def)
}
