package render_test

import (
	"encoding/json"
	"encoding/xml"
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

func genericRenderTest(mw noodle.Middleware, data interface{}) (w *httptest.ResponseRecorder, err error) {
	h := noodle.New(mw).Then(
		func(c context.Context, w http.ResponseWriter, r *http.Request) error {
			return render.Yield(c, 200, data)
		})

	r, _ := http.NewRequest("GET", "http://localhost/testId", nil)
	w = httptest.NewRecorder()
	err = h(context.TODO(), w, r)
	return
}

func TestTemplate(t *testing.T) {
	is := is.New(t)
	testData := TestStruct{2, "Hehehehe"}
	testTpl, _ := template.New("index").Parse("<b>{{ .A }}</b><i>{{ .B }}</i>")

	w, err := genericRenderTest(render.Template(testTpl), testData)
	is.NotErr(err)
	is.Equal(w.Header().Get("Content-Type"), "text/html")
	is.Equal(w.Body.String(), "<b>2</b><i>Hehehehe</i>")
}

func TestXML(t *testing.T) {
	is := is.New(t)
	testData := TestStruct{3, "Hahahahah"}

	w, err := genericRenderTest(render.XML, testData)
	is.NotErr(err)
	is.Equal(w.Header().Get("Content-Type"), "application/xml")

	var res TestStruct
	is.NotErr(xml.Unmarshal(w.Body.Bytes(), &res))
	is.Equal(res, testData)
}

func TestTextXML(t *testing.T) {
	is := is.New(t)
	testData := TestStruct{3, "Hahahahah"}

	w, err := genericRenderTest(render.TextXML, testData)
	is.NotErr(err)
	is.Equal(w.Header().Get("Content-Type"), "text/xml")

	var res TestStruct
	is.NotErr(xml.Unmarshal(w.Body.Bytes(), &res))
	is.Equal(res, testData)
}

func TestJSON(t *testing.T) {
	is := is.New(t)
	testData := TestStruct{3, "Hahahahah"}

	w, err := genericRenderTest(render.JSON, testData)
	is.NotErr(err)
	is.Equal(w.Header().Get("Content-Type"), "application/json")

	var res TestStruct
	is.NotErr(json.Unmarshal(w.Body.Bytes(), &res))
	is.Equal(res, testData)
}
