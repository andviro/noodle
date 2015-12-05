package main

import (
	"fmt"
	hr "github.com/andviro/noodle/adapt/httprouter"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/net/context"
	"net/http"
)

func index(c context.Context, w http.ResponseWriter, r *http.Request) error {
	fmt.Fprintln(w, "index")
	return nil
}

func products(c context.Context, w http.ResponseWriter, r *http.Request) error {
	params := hr.GetParams(c)
	fmt.Fprintf(w, "products: %s", params.ByName("id"))
	return nil
}

func main() {
	r := httprouter.New()
	n := hr.Default()
	r.GET("/", n.Then(index))
	r.GET("/products/:id", n.Then(products))
	http.ListenAndServe(":8080", r)
}
