package main

import (
	"fmt"
	"github.com/andviro/noodle/adapt/gorilla"
	mw "github.com/andviro/noodle/middleware"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "index")
}

func products(w http.ResponseWriter, r *http.Request) {
	vars := gorilla.GetVars(r)
	fmt.Fprintf(w, "products: %s", vars["id"])
}

func main() {
	r := mux.NewRouter()
	n := mw.Default(gorilla.Vars)
	r.Handle("/", n.Then(index))
	r.Handle("/products/{id}", n.Then(products))
	log.Fatal(http.ListenAndServe(":8080", r))
}
