package bind_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"reflect"

	"gopkg.in/andviro/noodle.v2"
	"gopkg.in/andviro/noodle.v2/bind"
)

type bindTestCase struct {
	Payload    string
	Value      interface{}
	Error      bool
	Middleware func(interface{}) noodle.Middleware
}

type testStruct struct {
	A int    `json:"a" form:"a"`
	B string `json:"b" form:"b"`
}

var bindTestCases = []bindTestCase{
	{"alskdjasdklj", testStruct{}, true, bind.Form},
	{"a=1&b=Ololo", testStruct{1, "Ololo"}, false, bind.Form},
	{`{"a": 1, "b": "Ololo"}`, testStruct{1, "Ololo"}, false, bind.JSON},
	{"{}", testStruct{}, false, bind.JSON},
	{"", testStruct{}, true, bind.JSON},
}

func TestBind(t *testing.T) {
	for _, tc := range bindTestCases {
		w := httptest.NewRecorder()
		r, _ := http.NewRequest("POST", "", bytes.NewBuffer([]byte(tc.Payload)))
		n := tc.Middleware(tc.Value)(func(w http.ResponseWriter, r *http.Request) {
			data, err := bind.Get(r)
			if err != nil && !tc.Error {
				t.Errorf("unexpected error: %+v", err)
			}
			if !reflect.DeepEqual(reflect.ValueOf(data).Elem().Interface(), tc.Value) {
				t.Errorf("expected %+v, got %+v", tc.Value, data)
			}
		})
		n(w, r)
	}
}

func TestBindPanicsOnPointer(t *testing.T) {
	defer func() {
		msg, ok := recover().(string)
		if !ok {
			t.Fatalf("should have thrown a string")
		}
		if msg != "bind to pointer is not allowed" {
			t.Error("invalid error value")
		}
	}()
	_ = bind.JSON(&bindTestCase{})
}
