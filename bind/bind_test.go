package bind_test

import (
	"bytes"
	"github.com/andviro/noodle"
	"github.com/andviro/noodle/bind"
	"golang.org/x/net/context"
	"gopkg.in/tylerb/is.v1"
	"net/http"
	"net/http/httptest"
	"testing"
)

type TestStruct struct {
	A int    `json:"a"`
	B string `json:"b"`
}

func bindHandlerFactory(is *is.Is) noodle.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
		s, ok := bind.GetData(ctx).(*TestStruct)
		is.True(ok)
		is.Equal(s.A, 1)
		is.Equal(s.B, "Ololo")
		return nil
	}
}

func TestBindJSON(t *testing.T) {
	is := is.New(t)
	buf := bytes.NewBuffer([]byte(`{"a": 1, "b": "Ololo"}`))
	emptyBuf := bytes.NewBuffer([]byte{})

	n := noodle.New(bind.JSON(TestStruct{})).Then(bindHandlerFactory(is))
	r, _ := http.NewRequest("POST", "http://localhost", buf)
	is.NotErr(n(context.TODO(), httptest.NewRecorder(), r))

	r, _ = http.NewRequest("POST", "http://localhost", emptyBuf)
	is.Err(n(context.TODO(), httptest.NewRecorder(), r))
}

func TestBindJSONPanicsOnPointer(t *testing.T) {
	is := is.New(t)
	var err interface{}

	func() {
		defer func() {
			err = recover()
		}()
		_ = bind.JSON(&TestStruct{})
	}()
	is.Equal(err.(string), "Bind to pointer is not allowed")
}
