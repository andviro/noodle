package render_test

import (
	"encoding/json"
	"github.com/andviro/noodle"
	"github.com/andviro/noodle/render"
	"golang.org/x/net/context"
	"gopkg.in/tylerb/is.v1"
	"html/template"
	"net/http"
	"net/http/httptest"
	"testing"
)

type TestStruct struct {
	A int    `json:"a"`
	B string `json:"b"`
}

func TestJSON(t *testing.T) {
	is := is.New(t)
	testData := TestStruct{1, "Ohohoho"}

	h := noodle.New(render.JSON).Then(
		func(c context.Context, w http.ResponseWriter, r *http.Request) error {
			return render.Yield(c, 200, &testData)
		})

	r, _ := http.NewRequest("GET", "http://localhost/testId", nil)
	w := httptest.NewRecorder()
	is.NotErr(h(context.TODO(), w, r))
	is.Equal(w.Header().Get("Content-Type"), "application/json")

	var res TestStruct
	is.NotErr(json.Unmarshal(w.Body.Bytes(), &res))
	is.Equal(res, testData)

}

func TestTemplate(t *testing.T) {
	is := is.New(t)
	testData := TestStruct{2, "Hehehehe"}
	testTpl, _ := template.New("index").Parse("<b>{{ .A }}</b><i>{{ .B }}</i>")

	h := noodle.New(render.Template(testTpl)).Then(
		func(c context.Context, w http.ResponseWriter, r *http.Request) error {
			return render.Yield(c, 200, &testData)
		})

	r, _ := http.NewRequest("GET", "http://localhost/testId", nil)
	w := httptest.NewRecorder()
	is.NotErr(h(context.TODO(), w, r))
	is.Equal(w.Header().Get("Content-Type"), "text/html")
	is.Equal(w.Body.String(), "<b>2</b><i>Hehehehe</i>")
}
