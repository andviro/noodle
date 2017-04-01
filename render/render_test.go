package render_test

import (
	"encoding/json"
	"encoding/xml"
	"gopkg.in/andviro/noodle.v2"
	"gopkg.in/andviro/noodle.v2/render"
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

func genericRenderTest(mw noodle.Middleware, data interface{}, mutators ...headerMutator) (w *httptest.ResponseRecorder) {
	h := noodle.New(mw).Then(
		func(w http.ResponseWriter, r *http.Request) {
			render.Yield(r, 200, data)
		})

	r, _ := http.NewRequest("GET", "http://localhost/testId", nil)
	for _, m := range mutators {
		m(r)
	}
	w = httptest.NewRecorder()
	h(w, r)
	return
}

func TestTemplate(t *testing.T) {
	is := is.New(t)
	testData := TestStruct{2, "Hehehehe"}
	testTpl, _ := template.New("index").Parse("<b>{{ .A }}</b><i>{{ .B }}</i>")

	w := genericRenderTest(render.Template(testTpl), testData)
	is.Equal(w.Header().Get("Content-Type"), "text/html;charset=utf-8")
	is.Equal(w.Body.String(), "<b>2</b><i>Hehehehe</i>")
}

func TestXML(t *testing.T) {
	is := is.New(t)
	testData := TestStruct{3, "Hahahahah"}

	w := genericRenderTest(render.XML, testData)
	is.Equal(w.Header().Get("Content-Type"), "application/xml;charset=utf-8")

	var res TestStruct
	is.NotErr(xml.Unmarshal(w.Body.Bytes(), &res))
	is.Equal(res, testData)
}

func TestTextXML(t *testing.T) {
	is := is.New(t)
	testData := TestStruct{3, "Hahahahah"}

	w := genericRenderTest(render.TextXML, testData)
	is.Equal(w.Header().Get("Content-Type"), "text/xml;charset=utf-8")

	var res TestStruct
	is.NotErr(xml.Unmarshal(w.Body.Bytes(), &res))
	is.Equal(res, testData)
}

func TestJSON(t *testing.T) {
	is := is.New(t)
	testData := TestStruct{3, "Hahahahah"}

	w := genericRenderTest(render.JSON, testData)
	is.Equal(w.Header().Get("Content-Type"), "application/json;charset=utf-8")

	var res TestStruct
	is.NotErr(json.Unmarshal(w.Body.Bytes(), &res))
	is.Equal(res, testData)
}

func TestContentTypeHTML(t *testing.T) {
	is := is.New(t)
	testData := TestStruct{2, "Hehehehe"}
	testTpl, _ := template.New("index").Parse("<b>{{ .A }}</b><i>{{ .B }}</i>")

	w := genericRenderTest(render.ContentType(testTpl), testData, func(r *http.Request) {
		r.Header.Set("Accept", "text/html")
	})
	is.Equal(w.Header().Get("Content-Type"), "text/html;charset=utf-8")
	is.Equal(w.Body.String(), "<b>2</b><i>Hehehehe</i>")
}

func TestContentTypeNil(t *testing.T) {
	is := is.New(t)
	testData := TestStruct{2, "Hehehehe"}

	w := genericRenderTest(render.ContentType(nil), testData, func(r *http.Request) {
		r.Header.Set("Accept", "text/html")
	})
	is.Equal(w.Header().Get("Content-Type"), "text/html;charset=utf-8")
	// XXX: too lazy to test this properly
}

func TestContentTypeJSON(t *testing.T) {
	is := is.New(t)
	testData := TestStruct{3, "Hahahahah"}

	w := genericRenderTest(render.ContentType(nil), testData, func(r *http.Request) {
		r.Header.Set("Accept", "application/json")
	})
	is.Equal(w.Header().Get("Content-Type"), "application/json;charset=utf-8")

	var res TestStruct
	is.NotErr(json.Unmarshal(w.Body.Bytes(), &res))
	is.Equal(res, testData)
}

func TestContentTypeXML(t *testing.T) {
	is := is.New(t)
	testData := TestStruct{3, "Hahahahah"}

	w := genericRenderTest(render.ContentType(nil), testData, func(r *http.Request) {
		r.Header.Set("Accept", "application/xml")
	})
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
		w := genericRenderTest(render.ContentType(nil), nil, func(r *http.Request) {
			if ct.Expected != "" {
				r.Header.Set("Accept", ct.Expected)
			}
		})
		is.Equal(w.Header().Get("Content-Type"), ct.Received)
	}
}
