package main

import (
	"fmt"
	"github.com/andviro/noodle/adapt/gorilla"
	mw "github.com/andviro/noodle/middleware"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
	"log"
	"net/http"
)

func index(c context.Context, w http.ResponseWriter, r *http.Request) error {
	fmt.Fprintln(w, "index")
	return nil
}

func products(c context.Context, w http.ResponseWriter, r *http.Request) error {
	vars := gorilla.GetVars(c)
	fmt.Fprintf(w, "products: %s", vars["id"])
	return nil
}

func main() {
	r := mux.NewRouter()
	n := mw.Default(gorilla.Vars)
	r.Handle("/", n.Then(index))
	r.Handle("/products/{id}", n.Then(products))
	log.Fatal(http.ListenAndServe(":8080", r))
}
