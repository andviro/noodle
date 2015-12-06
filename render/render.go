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
		return json.NewEncoder(w).Encode(res.data)
	}
}

// Yield puts arbitrary object into context for subsequent rendering into response
func Yield(c context.Context, data interface{}) {
	dest := c.Value(renderKey).(*renderResult)
	dest.mu.Lock() // better safe than sorry
	defer dest.mu.Unlock()
	dest.data = data
}
