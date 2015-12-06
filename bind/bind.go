package bind

import (
	"encoding/json"
	"github.com/andviro/noodle"
	"golang.org/x/net/context"
	"net/http"
	"reflect"
)

type key int

var (
	bindKey key = 0
)

// JSON constructs middleware that parses request body according to provided model
// and injects parsed object into context
func JSON(model interface{}) noodle.Middleware {
	typeModel := reflect.TypeOf(model)
	if typeModel.Kind() == reflect.Ptr {
		panic("Bind to pointer is not allowed")
	}
	return func(next noodle.Handler) noodle.Handler {
		return func(c context.Context, w http.ResponseWriter, r *http.Request) error {
			res := reflect.New(typeModel).Interface()
			err := json.NewDecoder(r.Body).Decode(res)
			if err != nil {
				return err
			}
			return next(context.WithValue(c, bindKey, res), w, r)
		}
	}
}

// GetData extracts data parsed from upstream Bind operation
func GetData(c context.Context) interface{} {
	return c.Value(bindKey)
}
