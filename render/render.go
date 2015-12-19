package render

import (
	"encoding/json"
	"github.com/andviro/noodle"
	"golang.org/x/net/context"
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

// serializerFunc is modelled after template's Execute method
type serializerFunc func(io.Writer, interface{}) error

// generic factory for a middleware that lifts handler result object from context
// and serializes it into HTTP ResponseWriter. Receives serializer function
func generic(s serializerFunc) noodle.Middleware {
	return func(next noodle.Handler) noodle.Handler {
		return func(c context.Context, w http.ResponseWriter, r *http.Request) error {
			var res renderResult

			err := next(context.WithValue(c, renderKey, &res), w, r)
			if err != nil {
				return err
			}
			w.Header().Set("Content-Type", "application/json")

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
var JSON = generic(func(w io.Writer, data interface{}) error {
	return json.NewEncoder(w).Encode(data)
})

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
