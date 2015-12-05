package render

import (
	"encoding/json"
	"github.com/andviro/noodle"
	"golang.org/x/net/context"
	"net/http"
)

var renderKey int = 0

type renderResult struct {
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
		return json.NewEncoder(w).Encode(res.data)
	}
}

// Yield puts arbitrary object into context for subsequent rendering into response
func Yield(c context.Context, data interface{}) {
	dest := c.Value(renderKey).(*renderResult)
	dest.data = data
}
