package render

import (
	"encoding/json"
	"github.com/andviro/noodle"
	"golang.org/x/net/context"
	"net/http"
	"sync"
)

var renderKey int = 0

type renderResult struct {
	mu   sync.RWMutex // guards data
	code int
	data interface{}
}

// JSON is a middleware that lifts handler result object from context
// and serializes it into HTTP ResponseWriter
func JSON(next noodle.Handler) noodle.Handler {
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
		return json.NewEncoder(w).Encode(res.data)
	}
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
