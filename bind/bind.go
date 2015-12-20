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

// Constructor is a generic function modelled after json.NewDecoder
type Constructor func(io.Reader) Decoder

// Decoder populates target object with data from request body
type Decoder interface {
	Decode(interface{}) error
}

func jsonC(r io.Reader) Decoder {
	return json.NewDecoder(r)
}

func formC(r io.Reader) Decoder {
	return form.NewDecoder(r)
}

// Generic is a middleware factory for request binding.
// Accepts Constructor and returns binder for model
func Generic(dc Constructor) func(interface{}) noodle.Middleware {
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
var JSON = Generic(jsonC)

// Form constructs middleware that parses request form according to provided model
// and injects parsed object into context
var Form = Generic(formC)

// GetData extracts data parsed from upstream Bind operation
func GetData(c context.Context) interface{} {
	return c.Value(bindKey)
}
