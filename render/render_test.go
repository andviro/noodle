package render_test

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"github.com/andviro/noodle"
	"github.com/andviro/noodle/render"
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

type headerMutator func(*http.Request)

func genericRenderTest(mw noodle.Middleware, data interface{}, mutators ...headerMutator) (w *httptest.ResponseRecorder, err error) {
	h := noodle.New(mw).Then(
		func(c context.Context, w http.ResponseWriter, r *http.Request) error {
			return render.Yield(c, 200, data)
		})

	r, _ := http.NewRequest("GET", "http://localhost/testId", nil)
	for _, m := range mutators {
		m(r)
	}
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
	is.Equal(w.Header().Get("Content-Type"), "text/html;charset=utf-8")
	is.Equal(w.Body.String(), "<b>2</b><i>Hehehehe</i>")
}

func TestXML(t *testing.T) {
	is := is.New(t)
	testData := TestStruct{3, "Hahahahah"}

	w, err := genericRenderTest(render.XML, testData)
	is.NotErr(err)
	is.Equal(w.Header().Get("Content-Type"), "application/xml;charset=utf-8")

	var res TestStruct
	is.NotErr(xml.Unmarshal(w.Body.Bytes(), &res))
	is.Equal(res, testData)
}

func TestTextXML(t *testing.T) {
	is := is.New(t)
	testData := TestStruct{3, "Hahahahah"}

	w, err := genericRenderTest(render.TextXML, testData)
	is.NotErr(err)
	is.Equal(w.Header().Get("Content-Type"), "text/xml;charset=utf-8")

	var res TestStruct
	is.NotErr(xml.Unmarshal(w.Body.Bytes(), &res))
	is.Equal(res, testData)
}

func TestJSON(t *testing.T) {
	is := is.New(t)
	testData := TestStruct{3, "Hahahahah"}

	w, err := genericRenderTest(render.JSON, testData)
	is.NotErr(err)
	is.Equal(w.Header().Get("Content-Type"), "application/json;charset=utf-8")

	var res TestStruct
	is.NotErr(json.Unmarshal(w.Body.Bytes(), &res))
	is.Equal(res, testData)
}

func TestContentTypeHTML(t *testing.T) {
	is := is.New(t)
	testData := TestStruct{2, "Hehehehe"}
	testTpl, _ := template.New("index").Parse("<b>{{ .A }}</b><i>{{ .B }}</i>")

	w, err := genericRenderTest(render.ContentType(testTpl), testData, func(r *http.Request) {
		r.Header.Set("Accept", "text/html")
	})
	is.NotErr(err)
	is.Equal(w.Header().Get("Content-Type"), "text/html;charset=utf-8")
	is.Equal(w.Body.String(), "<b>2</b><i>Hehehehe</i>")
}

func TestContentTypeNil(t *testing.T) {
	is := is.New(t)
	testData := TestStruct{2, "Hehehehe"}

	w, err := genericRenderTest(render.ContentType(nil), testData, func(r *http.Request) {
		r.Header.Set("Accept", "text/html")
	})
	is.NotErr(err)
	is.Equal(w.Header().Get("Content-Type"), "text/html;charset=utf-8")
	// XXX: too lazy to test this properly
}

func TestContentTypeJSON(t *testing.T) {
	is := is.New(t)
	testData := TestStruct{3, "Hahahahah"}

	w, err := genericRenderTest(render.ContentType(nil), testData, func(r *http.Request) {
		r.Header.Set("Accept", "application/json")
	})
	is.NotErr(err)
	is.Equal(w.Header().Get("Content-Type"), "application/json;charset=utf-8")

	var res TestStruct
	is.NotErr(json.Unmarshal(w.Body.Bytes(), &res))
	is.Equal(res, testData)
}

func TestContentTypeXML(t *testing.T) {
	is := is.New(t)
	testData := TestStruct{3, "Hahahahah"}

	w, err := genericRenderTest(render.ContentType(nil), testData, func(r *http.Request) {
		r.Header.Set("Accept", "application/xml")
	})
	is.NotErr(err)
	is.Equal(w.Header().Get("Content-Type"), "application/xml;charset=utf-8")

	var res TestStruct
	is.NotErr(xml.Unmarshal(w.Body.Bytes(), &res))
	is.Equal(res, testData)
}

func TestContentTypeParsesAccept(t *testing.T) {
	is := is.New(t)
	contentTypes := []struct {
		Expected string
		Received string
	}{
		{"application/xml", "application/xml;charset=utf-8"},
		{"application/json", "application/json;charset=utf-8"},
		{"text/xml", "text/xml;charset=utf-8"},
		{"text/html", "text/html;charset=utf-8"},
		{"", "application/json;charset=utf-8"},
	}

	for _, ct := range contentTypes {
		w, err := genericRenderTest(render.ContentType(nil), nil, func(r *http.Request) {
			if ct.Expected != "" {
				r.Header.Set("Accept", ct.Expected)
			}
		})
		is.NotErr(err)
		is.Equal(w.Header().Get("Content-Type"), ct.Received)
	}
}
