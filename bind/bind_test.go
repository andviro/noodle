package bind_test

import (
	"bytes"
	"github.com/ajg/form"
	"gopkg.in/andviro/noodle.v2"
	"gopkg.in/andviro/noodle.v2/bind"
	"gopkg.in/tylerb/is.v1"
	"net/http"
	"net/http/httptest"
	"testing"
)

type TestStruct struct {
	A int    `json:"a" form:"a"`
	B string `json:"b" form:"b"`
}

func bindHandlerFactory(is *is.Is) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s, ok := bind.GetData(r).(*TestStruct)
		is.True(ok)
		is.Equal(s.A, 1)
		is.Equal(s.B, "Ololo")
	}
}

func TestBindForm(t *testing.T) {
	is := is.New(t)
	testForm, _ := form.EncodeToString(TestStruct{1, "Ololo"})
	buf := bytes.NewBuffer([]byte(testForm))
	emptyBuf := bytes.NewBuffer([]byte("alskdjasdklj"))

	n := noodle.New(bind.Form(TestStruct{})).Then(bindHandlerFactory(is))
	r, _ := http.NewRequest("POST", "http://localhost", buf)
	n(httptest.NewRecorder(), r)

	r, _ = http.NewRequest("POST", "http://localhost", emptyBuf)
	n(httptest.NewRecorder(), r)
}

func TestBindJSON(t *testing.T) {
	is := is.New(t)
	buf := bytes.NewBuffer([]byte(`{"a": 1, "b": "Ololo"}`))
	emptyBuf := bytes.NewBuffer([]byte{})

	n := noodle.New(bind.JSON(TestStruct{})).Then(bindHandlerFactory(is))
	r, _ := http.NewRequest("POST", "http://localhost", buf)
	n(httptest.NewRecorder(), r)

	r, _ = http.NewRequest("POST", "http://localhost", emptyBuf)
	n(httptest.NewRecorder(), r)
}

func TestBindPanicsOnPointer(t *testing.T) {
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
