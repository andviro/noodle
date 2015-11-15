package main

import (
	"fmt"
	"github.com/andviro/noodle"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
	"log"
	"net/http"
)

// GorillaVars injects Gorilla mux route variables into context
func GorillaVars(next noodle.Handler) noodle.Handler {
	return func(c context.Context, w http.ResponseWriter, r *http.Request) error {
		withVars := context.WithValue(c, "Vars", mux.Vars(r))
		return next(withVars, w, r)
	}
}

func index(c context.Context, w http.ResponseWriter, r *http.Request) error {
	fmt.Fprintln(w, "index")
	return nil
}

func products(c context.Context, w http.ResponseWriter, r *http.Request) error {
	vars := c.Value("Vars").(map[string]string)
	fmt.Fprintf(w, "products: %s", vars["id"])
	return nil
}

func main() {
	r := mux.NewRouter()
	n := noodle.Default(GorillaVars)
	r.Handle("/", n.Then(index))
	r.Handle("/products/{id}", n.Then(products))
	log.Fatal(http.ListenAndServe(":8080", r))
}
