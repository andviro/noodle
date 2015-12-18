package bind

import (
	"encoding/json"
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

type DecoderConstructor interface {
	NewDecoder(io.Reader) Decoder
}

type Decoder interface {
	Decode(interface{}) error
}

// Generic
func Generic(dc DecoderConstructor) {
	return func(model interface{}) noodle.Middleware {
		typeModel := reflect.TypeOf(model)
		if typeModel.Kind() == reflect.Ptr {
			panic("Bind to pointer is not allowed")
		}
		return func(next noodle.Handler) noodle.Handler {
			return func(c context.Context, w http.ResponseWriter, r *http.Request) error {
				res := reflect.New(typeModel).Interface()
				err := dc.NewDecoder(r.Body).Decode(res)
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
var JSON = Generic(jsonConstructor)

// GetData extracts data parsed from upstream Bind operation
func GetData(c context.Context) interface{} {
	return c.Value(bindKey)
}
