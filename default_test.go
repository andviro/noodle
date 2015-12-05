package noodle_test

import (
	"github.com/andviro/noodle"
	"gopkg.in/tylerb/is.v1"
	"testing"
)

func TestDefault(t *testing.T) {
	is := is.New(t)
	def := noodle.Default()
	is.Equal(len(def), 2)
}
