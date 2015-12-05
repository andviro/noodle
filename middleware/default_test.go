package middleware_test

import (
	mw "github.com/andviro/noodle/middleware"
	"gopkg.in/tylerb/is.v1"
	"testing"
)

func TestDefault(t *testing.T) {
	is := is.New(t)
	def := mw.Default()
	is.Equal(len(def), 3)
}
