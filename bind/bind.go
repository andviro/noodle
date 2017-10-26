package bind

import (
	"encoding/json"
	"io"
	"net/http"
	"reflect"

	"github.com/ajg/form"
	"gopkg.in/andviro/noodle.v2"
)

type key int

const bindKey key = iota

type decoderResult struct {
	val interface{}
	err error
}

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
			panic("bind to pointer is not allowed")
		}
		return func(next http.HandlerFunc) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				var res decoderResult
				res.val = reflect.New(typeModel).Interface()
				res.err = dc(r.Body).Decode(res.val)
				next(w, noodle.WithValue(r, bindKey, res))
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

// GetData extracts data parsed from upstream Bind operation, discarding
// decoding error.  Deprecated in favor of Get() function.
func GetData(r *http.Request) (res interface{}) {
	res, _ = Get(r)
	return
}

// Get extracts data parsed from upstream Bind operation, along with the decode error.
func Get(r *http.Request) (val interface{}, err error) {
	res, ok := noodle.Value(r, bindKey).(decoderResult)
	if !ok {
		return
	}
	return res.val, res.err
}
