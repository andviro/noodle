package bind

import (
	"encoding/json"
	"github.com/ajg/form"
	"github.com/andviro/noodle"
	"golang.org/x/net/context"
	"io"
	"net/http"
	"reflect"
)

type key int

var (
	bindKey key = 0
)

type constructor func(io.Reader) decoder

type decoder interface {
	Decode(interface{}) error
}

func jsonC(r io.Reader) decoder {
	return json.NewDecoder(r)
}

func formC(r io.Reader) decoder {
	return form.NewDecoder(r)
}

// generic middleware factory for request binding.
// Accepts constructor type that receives io.Reader and returns decoder
func generic(dc constructor) func(interface{}) noodle.Middleware {
	return func(model interface{}) noodle.Middleware {
		typeModel := reflect.TypeOf(model)
		if typeModel.Kind() == reflect.Ptr {
			panic("Bind to pointer is not allowed")
		}
		return func(next noodle.Handler) noodle.Handler {
			return func(c context.Context, w http.ResponseWriter, r *http.Request) error {
				res := reflect.New(typeModel).Interface()
				err := dc(r.Body).Decode(res)
				if err != nil {
					return err
				}
				return next(context.WithValue(c, bindKey, res), w, r)
			}
		}
	}
}

// JSON constructs middleware that parses request body according to provided model
// and injects parsed object into context
var JSON = generic(jsonC)

// Form constructs middleware that parses request form according to provided model
// and injects parsed object into context
var Form = generic(formC)

// GetData extracts data parsed from upstream Bind operation
func GetData(c context.Context) interface{} {
	return c.Value(bindKey)
}
