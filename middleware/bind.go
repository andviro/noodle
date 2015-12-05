package middleware

import (
	"encoding/json"
	"github.com/andviro/noodle"
	"golang.org/x/net/context"
	"net/http"
	"reflect"
)

// BindJSON constructs middleware that parses request body according to provided model
// and injects parsed object into context
func BindJSON(model interface{}) noodle.Middleware {
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

// GetBindData extracts data parsed from upstream Bind operation
func GetBindData(c context.Context) interface{} {
	return c.Value(bindKey)
}
