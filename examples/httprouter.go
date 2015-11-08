package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/context"
	"net/http"
	"noodle"
)

// HR adapts noodle middleware chain to httprouter handler signature.
// httprouter.Params are injected into context under key "Params"
func HR(h noodle.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		ctx := context.WithValue(context.TODO(), "Params", p)
		h(ctx, w, r)
	}
}

func index(c context.Context, w http.ResponseWriter, r *http.Request) error {
	fmt.Fprintln(w, "index")
	return nil
}

func products(c context.Context, w http.ResponseWriter, r *http.Request) error {
	params := c.Value("Params").(httprouter.Params)
	fmt.Fprintf(w, "products: %s", params.ByName("id"))
	return nil
}

func main() {
	r := httprouter.New()
	n := noodle.Default()
	r.GET("/", HR(n.Then(index)))
	r.GET("/products/:id", HR(n.Then(products)))
	http.ListenAndServe(":8080", r)
}
