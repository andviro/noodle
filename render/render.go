package render

import (
	"encoding/json"
	"encoding/xml"
	"github.com/andviro/noodle"
	"golang.org/x/net/context"
	"html/template"
	"io"
	"net/http"
	"sync"
)

var renderKey int = 0

type renderResult struct {
	mu   sync.RWMutex // guards data
	code int
	data interface{}
}

// SerializerFunc is modelled after template's Execute method.
// Used by Generic middleware factory to create specific rendering middlewares
type SerializerFunc func(io.Writer, interface{}) error

// Generic factory for a middleware that lifts handler's data object from context
// and serializes it into HTTP ResponseWriter. Receives SerializerFunc and content type
func Generic(s SerializerFunc, contentType string) noodle.Middleware {
	return func(next noodle.Handler) noodle.Handler {
		return func(c context.Context, w http.ResponseWriter, r *http.Request) error {
			var res renderResult

			err := next(context.WithValue(c, renderKey, &res), w, r)
			if err != nil {
				return err
			}
			w.Header().Set("Content-Type", contentType)

			res.mu.RLock()
			defer res.mu.RUnlock()
			if res.code != 0 {
				w.WriteHeader(res.code)
			}
			return s(w, res.data)
		}
	}
}

// JSON serializes result object into JSON format
var JSON = Generic(func(w io.Writer, data interface{}) error {
	return json.NewEncoder(w).Encode(data)
}, "application/json")

// XML serializes result object into "application/xml" content type. Use TextXML for "text/xml" output.
var XML = Generic(func(w io.Writer, data interface{}) error {
	return xml.NewEncoder(w).Encode(data)
}, "application/xml")

// TextXML is the same as XML, but with "text/xml" content type
var TextXML = Generic(func(w io.Writer, data interface{}) error {
	return xml.NewEncoder(w).Encode(data)
}, "text/xml")

// Template creates middleware that applies pre-compiled template to handler's data object
func Template(tpl *template.Template) noodle.Middleware {
	return Generic(tpl.Execute, "text/html")
}

// Yield puts arbitrary data into context for subsequent rendering into response.
// The first argument of Yield is a HTTP status code.
func Yield(c context.Context, code int, data interface{}) error {
	dest := c.Value(renderKey).(*renderResult)
	dest.mu.Lock() // better safe than sorry
	defer dest.mu.Unlock()
	dest.code = code
	dest.data = data
	return nil
}
